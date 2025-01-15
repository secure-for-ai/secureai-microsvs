package sqlBuilderV3

import (
	"reflect"
	"strings"
	"sync"
	"unsafe"

	"github.com/secure-for-ai/secureai-microsvs/db"
	"github.com/secure-for-ai/secureai-microsvs/util"
)

type Type int

const (
	RawType Type = iota
	InsertType
	DeleteType
	UpdateType
	SelectType
	UpsertType
)

type Columns []string
type Map map[string]any

type Table interface {
	GetTableName() string
}

type fromItem interface {
	setAliasName(string)
	writeTo(*Writer)
	destroy()
}

var fromTablePool = sync.Pool{
	New: func() any {
		return new(fromTable)
	},
}

type fromTable struct {
	tableName string
	alias     string
}

func createFromTable(tableName string, alias string) fromItem {
	from := fromTablePool.Get().(*fromTable)
	from.tableName = tableName
	from.alias = alias
	return from
}

func (from *fromTable) setAliasName(name string) {
	from.alias = name
}

func (from *fromTable) writeTo(w *Writer) {
	w.WriteString(from.tableName)
	if len(from.alias) > 0 {
		w.WriteString(" AS ")
		w.WriteString(from.alias)
	}
}

func (from *fromTable) destroy() {
	fromTablePool.Put(from)
}

var fromStmtPool = sync.Pool{
	New: func() any {
		return new(fromStmt)
	},
}

type fromStmt struct {
	stmt  *Stmt
	alias string
}

func createFromStmt(stmt *Stmt, alias string) fromItem {
	from := fromStmtPool.Get().(*fromStmt)
	from.stmt = stmt
	from.alias = alias
	return from
}

func (from *fromStmt) setAliasName(name string) {
	from.alias = name
}

func (from *fromStmt) writeTo(w *Writer) {
	w.WriteByte('(')
	from.stmt.WriteTo(w)

	if len(from.alias) > 0 {
		w.WriteString(") AS ")
		w.WriteString(from.alias)
	} else {
		w.WriteByte(')')
	}
}

func (from *fromStmt) destroy() {
	from.stmt.Destroy()
	fromStmtPool.Put(from)
}

type Stmt struct {
	RefTable *Table

	tableInto string
	tableFrom []fromItem

	where Cond
	// tracker internal created conds. Reset() only destroy
	// refed conds. This can avoid double free.
	whereRef []Cond

	GroupByStr *stringWriter
	having     Cond
	// tracker internal created conds. Reset() only destroy
	// refed conds. This can avoid double free.
	havingRef  []Cond
	OrderByStr *stringWriter

	Offset int
	LimitN int

	InsertCols   []string
	InsertValues valExpr2DList

	SetCols *condExpr

	SelectCols []string

	sqlType Type
}

func Insert(data ...any) *Stmt {
	return SQL().Insert(data...)
}

func InsertOne(data any) *Stmt {
	return SQL().InsertOne(data)
}

func InsertBulk(data any) *Stmt {
	return SQL().InsertBulk(data)
}

func Delete(data ...any) *Stmt {
	return SQL().Delete(data...)
}

func Update(data ...any) *Stmt {
	return SQL().Update(data...)
}

func Select(data ...any) *Stmt {
	return SQL().Select(data...)
}

var stmtPool = sync.Pool{
	New: func() any {
		stmt := &Stmt{}
		stmt.Init()
		return stmt
	},
}

func SQL() *Stmt {
	return stmtPool.Get().(*Stmt)
}

// Init reset all the statement's fields
func (stmt *Stmt) Init() {
	stmt.RefTable = nil

	stmt.tableInto = ""
	stmt.tableFrom = make([]fromItem, 0, 2)

	stmt.where = condEmpty{}
	stmt.whereRef = make([]Cond, 0, 2)
	stmt.GroupByStr = new(stringWriter)
	stmt.having = condEmpty{}
	stmt.havingRef = make([]Cond, 0, 2)
	stmt.OrderByStr = new(stringWriter)

	stmt.Offset = 0
	stmt.LimitN = 0

	stmt.InsertCols = []string{}
	stmt.InsertValues = newValExpr2DList(2)
	stmt.SetCols = Expr("")
	stmt.SelectCols = []string{}

	stmt.sqlType = RawType
}

func (stmt *Stmt) Reset() {
	stmt.RefTable = nil

	stmt.tableInto = ""
	for _, from := range stmt.tableFrom {
		from.destroy()
	}
	stmt.tableFrom = stmt.tableFrom[:0]

	stmt.where = condEmpty{}
	for _, cond := range stmt.whereRef {
		cond.Destroy()
	}
	stmt.whereRef = stmt.whereRef[:0]
	stmt.GroupByStr.Reset()
	stmt.having = condEmpty{}
	for _, cond := range stmt.havingRef {
		cond.Destroy()
	}
	stmt.havingRef = stmt.havingRef[:0]
	stmt.OrderByStr.Reset()

	stmt.Offset = 0
	stmt.LimitN = 0

	stmt.InsertCols = stmt.InsertCols[:0]
	stmt.InsertValues.reset()
	stmt.SetCols.Reset()
	stmt.SelectCols = stmt.SelectCols[:0]

	stmt.sqlType = RawType
}

func (stmt *Stmt) Destroy() {
	stmt.Reset()
	stmtPool.Put(stmt)
}

// Todo replace with the fast version of sync map
var structColumnCache = sync.Map{}

type stringWriter struct {
	strings.Builder
}

type builderInternalType struct {
	addr *builderInternalType
	buf  []byte
}

func (w *stringWriter) Bytes() []byte {
	b := (*builderInternalType)(unsafe.Pointer(w))
	return b.buf
}

func (w *stringWriter) Reset() {
	b := (*builderInternalType)(unsafe.Pointer(w))
	b.buf = b.buf[:0]
}

func (w *stringWriter) Destroy() {
	w.Reset()
	bufPool.Put(w)
}

var bufPool = sync.Pool{
	New: func() any {
		return &stringWriter{}
	},
}

func buildColumnsInternal(v reflect.Value, vType reflect.Type) []string {
	// construct the unique name of the struct with zero-copy method
	w := bufPool.Get().(*stringWriter)
	w.WriteString(vType.PkgPath())
	w.WriteByte('.')
	w.WriteString(vType.Name())
	structFullName := w.String()
	structColNames, ok := structColumnCache.Load(structFullName)

	// write the column name of the struct from the cache
	if ok {
		// free *stringWriter
		w.Destroy()
		return structColNames.([]string)
	}

	// avoid extend the slice cap which causes memory reallocation
	numField := v.NumField()
	tmpColNames := make([]string, numField)
	for i, il := 0, numField; i < il; i++ {
		// Get column name, tag start with "pg" or the field Name
		var colName string
		fieldInfo := vType.Field(i)
		if colName = fieldInfo.Tag.Get(db.Tag); colName == "" {
			colName = fieldInfo.Name
		}
		tmpColNames[i] = colName
	}

	// store the column name of the struct into the cache
	structColumnCache.Store(strings.Clone(structFullName), tmpColNames)
	// free *stringWriter
	w.Destroy()

	return tmpColNames

}
func buildColumns(colNames *[]string, column any) {
	v := util.ReflectValue(column)
	vType := v.Type()

	if vType.Kind() == reflect.Struct {
		*colNames = append(*colNames, buildColumnsInternal(v, vType)...)
	}
}

//go:linkname valueInterface reflect.valueInterface
func valueInterface(v reflect.Value, safe bool) any

func buildValues(curData any) *valExprList {
	v := util.ReflectValue(curData)
	vType := v.Type()
	if vType.Kind() == reflect.Struct {
		numField := v.NumField()
		values := getValExprListWithSize(numField) //make([]condExpr, numField)
		for i, il := 0, numField; i < il; i++ {
			// Get value
			fieldValue := v.Field(i)
			values.SetIth(i, db.Para, valueInterface(fieldValue, false))
			// switch fieldValue.Kind() {
			// default:
			// 	values.SetIth(i, db.Para, valueInterface(fieldValue, false))
			// case reflect.Bool:
			// 	values.SetIth(i, db.Para, fieldValue.Bool())
			// case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			// 	// args := getArgsWithSize(1)
			// 	// *args = append(*args, fieldValue.Int())
			// 	values.SetIth(i, db.Para, fieldValue.Int())
			// 	// argsPool.Put(args)
			// case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			// 	values.SetIth(i, db.Para, fieldValue.Uint())
			// case reflect.Float32, reflect.Float64:
			// 	values.SetIth(i, db.Para, fieldValue.Float())
			// case reflect.Complex64, reflect.Complex128:
			// 	values.SetIth(i, db.Para, fieldValue.Complex())
			// case reflect.String:
			// 	values.SetIth(i, db.Para, fieldValue.String())
			// }

			// (*values)[i].sql = db.Para
			// (*values)[i].args[0] = val

		}
		return values
	}
	return nil
}

// Into sets insert table name
func (stmt *Stmt) IntoTable(table any) *Stmt {
	switch table := table.(type) {
	case Table:
		stmt.tableInto = table.GetTableName()
	case string:
		stmt.tableInto = table
	}
	return stmt
}

func (stmt *Stmt) IntoColumns(column any, cols ...string) *Stmt {
	switch column := column.(type) {
	case []string:
		stmt.InsertCols = append(stmt.InsertCols, column...)
	case Columns:
		stmt.InsertCols = append(stmt.InsertCols, column...)
	case string:
		stmt.InsertCols = append(stmt.InsertCols, column)
		stmt.InsertCols = append(stmt.InsertCols, cols...)
	case Map:
		stmt.buildInsertColsByMap(column)
	default:
		buildColumns(&stmt.InsertCols, column)
	}
	return stmt
}

// Values store the insertion data, optimized for one record and support bulk insertion as well
func (stmt *Stmt) Values(data ...any) *Stmt {
	switch len(data) {
	case 0:
		return stmt
	case 1:
		stmt.ValuesOne(data[0])
	default:
		return stmt.ValuesBulk(data)
	}

	return stmt
}

func (stmt *Stmt) ValuesOne(data any) *Stmt {
	if data == nil {
		return stmt
	}

	stmt.InsertValues.grow(1)
	curData := data
	switch curData := curData.(type) {
	case []any:
		InsertValues := getValExprListWithSize(len(curData))
		for i, val := range curData {
			if e, ok := val.(*condExpr); ok {
				InsertValues.SetIthWithExpr(i, e)
			} else {
				InsertValues.SetIth(i, db.Para, val)
			}
		}
		stmt.InsertValues.append(InsertValues)
	case Map:
		stmt.valuesOneMap(curData)
	case *Stmt:
		stmt.From(curData)
	default:
		if len(stmt.InsertCols) == 0 {
			buildColumns(&stmt.InsertCols, curData)
		}
		insertValues := buildValues(curData)
		if insertValues != nil {
			stmt.InsertValues.append(insertValues)
		}
	}

	return stmt
}

func (stmt *Stmt) ValuesBulk(data any) *Stmt {
	dataR := util.ReflectValue(data)
	dataType := dataR.Kind()
	if dataType != reflect.Slice && dataType != reflect.Array {
		return stmt
	}
	return stmt.valuesBulkInternal(&dataR)
}

func (stmt *Stmt) valuesBulkInternal(data *reflect.Value) *Stmt {
	dataLen := data.Len()
	if dataLen == 0 {
		return stmt
	}

	// update the insert columns
	data0 := data.Index(0).Interface()
	switch data0 := data0.(type) {
	case Map:
		stmt.buildInsertColsByMap(data0)
	default:
		buildColumns(&stmt.InsertCols, data0)
	}

	// loading the data
	stmt.InsertValues.grow(dataLen)
	for i := 0; i < dataLen; i++ {
		curData := data.Index(i).Interface()
		switch curData := curData.(type) {
		case []any:
			insertValues := getValExprListWithSize(len(curData))
			for j, val := range curData {
				if e, ok := val.(*condExpr); ok {
					insertValues.SetIthWithExpr(j, e)
				} else {
					insertValues.SetIth(j, db.Para, val)
				}
			}
			stmt.InsertValues.append(insertValues)
		case Map:
			insertValues := getValExprListWithSize(len(curData))
			for j, col := range stmt.InsertCols {
				val := curData[col]
				if e, ok := val.(*condExpr); ok {
					insertValues.SetIthWithExpr(j, e)
				} else {
					insertValues.SetIth(j, db.Para, val)
				}
			}
			stmt.InsertValues.append(insertValues)
		default:
			insertValues := buildValues(curData)
			if insertValues != nil {
				stmt.InsertValues.append(insertValues)
			}
		}
	}

	return stmt
}

func (stmt *Stmt) SelectColumns(column any, cols ...string) *Stmt {
	switch column := column.(type) {
	case []string:
		stmt.SelectCols = append(stmt.SelectCols, column...)
	case Columns:
		stmt.SelectCols = append(stmt.SelectCols, column...)
	case string:
		stmt.SelectCols = append(stmt.SelectCols, column)
		stmt.SelectCols = append(stmt.SelectCols, cols...)
	default:
		buildColumns(&stmt.SelectCols, column)
	}
	return stmt
}

// From sets from subject(can be a table name in string or a builder pointer) and its alias
func (stmt *Stmt) From(subject any, alias ...string) *Stmt {
	var from fromItem
	switch subject := subject.(type) {
	case *Stmt:
		//subquery should be a select statement, and we only accept one select stmt
		from = createFromStmt(subject, "")
	case Table:
		from = createFromTable(subject.GetTableName(), "")
	case string:
		from = createFromTable(subject, "")
	default:
		return stmt
	}

	if len(alias) > 0 {
		from.setAliasName(alias[0])
	}

	stmt.tableFrom = append(stmt.tableFrom, from)
	return stmt
}

// Insert SQL
func (stmt *Stmt) Insert(data ...any) *Stmt {
	switch len(data) {
	case 0:
		break
	default:
		// if data is a double array, you need to call IntoColumns afterwards.
		// Otherwise the order of the values should be the same order of the columns in the table.
		// Support Bulk Insertion
		stmt.Values(data...)
		stmt.IntoTable(data[0])
	}
	if stmt.sqlType == RawType {
		stmt.sqlType = InsertType
	}
	return stmt
}

func (stmt *Stmt) InsertOne(data any) *Stmt {
	if stmt.sqlType == RawType {
		stmt.sqlType = InsertType
	}
	stmt.ValuesOne(data)
	stmt.IntoTable(data)
	return stmt
}

// Insert SQL
func (stmt *Stmt) InsertBulk(data any) *Stmt {
	if stmt.sqlType == RawType {
		stmt.sqlType = InsertType
	}

	dataR := util.ReflectValue(data)
	dataType := dataR.Kind()
	if dataType != reflect.Slice && dataType != reflect.Array {
		return stmt
	}

	switch dataR.Len() {
	case 0:
		break
	default:
		// if data is a double array, you need to call IntoColumns afterwards.
		// Otherwise the order of the values should be the same order of the columns in the table.
		// Support Bulk Insertion
		stmt.valuesBulkInternal(&dataR)
		stmt.IntoTable(dataR.Index(0).Interface())
	}

	return stmt
}

// Delete sets delete SQL
func (stmt *Stmt) Delete(data ...any) *Stmt {
	l := len(data)
	if l >= 1 {
		stmt.From(data[0])
	}
	if l >= 2 {
		stmt.And(data[1], data[2:]...)
	}
	if stmt.sqlType == RawType {
		stmt.sqlType = DeleteType
	}
	return stmt
}

// Update
func (stmt *Stmt) Update(data ...any) *Stmt {
	l := len(data)
	if l >= 1 {
		stmt.Set(data[0])
		stmt.From(data[0])
	}
	if l >= 2 {
		stmt.And(data[1], data[2:]...)
	}
	if stmt.sqlType == RawType {
		stmt.sqlType = UpdateType
	}
	return stmt
}

// Select SQL
func (stmt *Stmt) Select(data ...any) *Stmt {
	l := len(data)
	if l >= 1 {
		if _, ok := data[0].(string); !ok {
			stmt.SelectColumns(data[0])
		}
		stmt.From(data[0])
	}
	if l >= 2 {
		stmt.And(data[1], data[2:]...)
	}
	if stmt.sqlType == RawType {
		stmt.sqlType = SelectType
	}
	return stmt
}

// Incr Generate  "Update ... Set column = column + arg" statement
func (stmt *Stmt) Incr(col string, args ...any) *Stmt {
	stmt.SetCols.appendInc(col, args...)
	return stmt
}

// Decr Generate  "Update ... Set column = column - arg" statement
func (stmt *Stmt) Decr(col string, args ...any) *Stmt {
	stmt.SetCols.appendDec(col, args...)
	return stmt
}

// setExpr Generate  "Update ... Set column = {expr}" statement
// if you want to use writeTo internal builtin functions without parameters like NOW(),
// then you'd better to call Set(col, Expr("Now()"))
// Todo support expr as SQLStmt
func (stmt *Stmt) setExpr(col string, expr any, args ...any) *Stmt {
	switch e := expr.(type) {
	case string:
		if len(args) > 0 {
			// set("col", "col||??", "test") => writeTo: col = col||??, args: "test"
			stmt.SetCols.appendSet(col, e, args...)
			return stmt
		}
	case *condExpr:
		stmt.SetCols.appendSet(col, e.String(), e.args...)
		return stmt
	}
	stmt.SetCols.appendEq(col, expr)
	return stmt
}

// setMap Generate  "Update ... Set col1 = {expr1}, col1 = {expr2}" statement
// {"username": "bob", "age": 10, "createTime": Expr("Now()"} =>
// SQL: username = ?? , age = ??, createTime = NOW()
// Args: ["bob", 10]
// Todo support expr as SQLStmt
func (stmt *Stmt) setMap(exprs Map) *Stmt {
	// avoid extend the slice cap which causes memory reallocation
	for col, val := range exprs {
		if e, ok := val.(*condExpr); ok {
			stmt.SetCols.appendSet(col, e.String(), e.args...)
		} else {
			stmt.SetCols.appendEq(col, val)
		}
	}
	return stmt
}

func (stmt *Stmt) setStruct(data any) *Stmt {
	// check whether data is struct
	// reflect the exact value of the data regardless of whether it's a ptr or struct
	v := util.ReflectValue(data)
	vType := v.Type()
	if vType.Kind() == reflect.Struct {

		var colNames []string = buildColumnsInternal(v, vType)
		numField := v.NumField()
		// avoid extend the slice cap which causes memory reallocation
		for i, il := 0, numField; i < il; i++ {
			// Get column name, tag start with "pg" or the field Name
			colName := colNames[i]

			// Get value
			fieldValue := v.Field(i)
			stmt.SetCols.appendEq(colName, valueInterface(fieldValue, false))
			// switch fieldValue.Kind() {
			// default:
			// 	stmt.SetCols.addParam(colName, Expr(db.Para, valueInterface(fieldValue, false)))
			// case reflect.Bool:
			// 	stmt.SetCols.addParam(colName, Expr(db.Para, fieldValue.Bool()))
			// case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			// 	stmt.SetCols.addParam(colName, Expr(db.Para, fieldValue.Int()))
			// case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			// 	stmt.SetCols.addParam(colName, Expr(db.Para, fieldValue.Uint()))
			// case reflect.Float32, reflect.Float64:
			// 	stmt.SetCols.addParam(colName, Expr(db.Para, fieldValue.Float()))
			// case reflect.Complex64, reflect.Complex128:
			// 	stmt.SetCols.addParam(colName, Expr(db.Para, fieldValue.Complex()))
			// case reflect.String:
			// 	stmt.SetCols.addParam(colName, Expr(db.Para, fieldValue.String()))
			// }
		}
	}
	return stmt
}

func (stmt *Stmt) Set(data any, args ...any) *Stmt {
	switch data := data.(type) {
	case string:
		argLen := len(args)
		if argLen >= 1 {
			stmt.setExpr(data, args[0], args[1:]...)
		}
	case Map:
		stmt.setMap(data)
	case *condExpr:
		stmt.SetCols.appendExpr(data)
	default:
		// assume the input is either a struct ptr or a struct
		stmt.setStruct(data)
	}
	return stmt
}

func (stmt *Stmt) Where(query any, args ...any) *Stmt {
	stmt.catCond(&stmt.where, &stmt.whereRef, And, query, args...)
	return stmt
}

// concat an existing Cond and a new Cond statement with Op
func (stmt *Stmt) catCond(c *Cond, ref *[]Cond, OpFunc func(cond ...Cond) Cond, query any, args ...any) {
	switch query := query.(type) {
	case string:
		// assume the input query is not empty
		cond := Expr(query, args...)
		// self created cond is stored in the ref
		*ref = append(*ref, cond)

		if (*c).IsValid() {
			// construct new cond with zero-copy
			conds := condAndPool.Get().(*condAnd)
			*conds = append(*conds, *c, cond)
			*c = OpFunc(*conds...)
			conds.Destroy()
			// new *c is either And or Or cond
			*ref = append(*ref, *c)
		} else {
			// *c is condEmpty{}, so replace it with new created cond
			*c = cond
		}
	case Map:
		conds := condAndPool.Get().(*condAnd)
		if _, ok := (*c).(condEmpty); !ok {
			*conds = append(*conds, *c)
		}
		catCondMap(ref, query, conds)
		*c = OpFunc(*conds...)
		if len(*conds) >= 2 {
			// we must construct a new cond either CondAnd or CondOr,
			// thus store *c into the ref
			*ref = append(*ref, *c)
		}
		conds.Destroy()
	case Cond:
		conds := condAndPool.Get().(*condAnd)
		if _, ok := (*c).(condEmpty); !ok {
			*conds = append(*conds, *c)
		}
		*conds = append(*conds, query)
		for _, v := range args {
			if vv, ok := v.(Cond); ok {
				*conds = append(*conds, vv)
			}
		}
		*c = OpFunc(*conds...)
		if len(*conds) >= 2 {
			// we must construct a new cond either CondAnd or CondOr,
			// thus store *c into the ref
			*ref = append(*ref, *c)
		}
		conds.Destroy()
	default:
		// TODO: not support condition type
	}
}

// And add Where & and statement
func (stmt *Stmt) And(query any, args ...any) *Stmt {
	stmt.catCond(&stmt.where, &stmt.whereRef, And, query, args...)
	return stmt
}

// Or add Where & Or statement
func (stmt *Stmt) Or(query any, args ...any) *Stmt {
	stmt.catCond(&stmt.where, &stmt.whereRef, Or, query, args...)
	return stmt
}

// GroupBy generate "Group By keys" statement
func (stmt *Stmt) GroupBy(keys ...string) *Stmt {
	if len(keys) == 0 {
		return stmt
	}

	groupByStr := stmt.GroupByStr

	if groupByStr.Len() > 0 {
		groupByStr.WriteString(", ")
	}

	bufferJoin(groupByStr, keys, ", ")
	return stmt
}

// GroupBy generate "Having conditions" statement
func (stmt *Stmt) Having(query any, args ...any) *Stmt {
	stmt.catCond(&stmt.having, &stmt.havingRef, And, query, args...)
	return stmt
}

// GroupBy generate "Having conditions" statement && conditions
func (stmt *Stmt) HavingAnd(query any, args ...any) *Stmt {
	stmt.catCond(&stmt.having, &stmt.havingRef, And, query, args...)
	return stmt
}

// GroupBy generate "Having conditions" statement || conditions
func (stmt *Stmt) HavingOr(query any, args ...any) *Stmt {
	stmt.catCond(&stmt.having, &stmt.havingRef, Or, query, args...)
	return stmt
}

func bufferJoin(w *stringWriter, elems []string, sep string) {
	w.WriteString(elems[0])
	for _, s := range elems[1:] {
		w.WriteString(sep)
		w.WriteString(s)
	}
}

// OrderBy generate "Order By order" statement
func (stmt *Stmt) OrderBy(order ...string) *Stmt {
	if len(order) == 0 {
		return stmt
	}

	orderByStr := stmt.OrderByStr

	if orderByStr.Len() > 0 {
		orderByStr.WriteString(", ")
	}

	bufferJoin(orderByStr, order, ", ")
	return stmt
}

// Desc generate `ORDER BY xx DESC`
func (stmt *Stmt) Desc(colNames ...string) *Stmt {
	if len(colNames) == 0 {
		return stmt
	}

	orderByStr := stmt.OrderByStr
	if orderByStr.Len() > 0 {
		orderByStr.WriteString(", ")
	}

	bufferJoin(orderByStr, colNames, " DESC, ")
	orderByStr.WriteString(" DESC")

	return stmt
}

// Asc generate `ORDER BY xx ASC`
func (stmt *Stmt) Asc(colNames ...string) *Stmt {
	if len(colNames) == 0 {
		return stmt
	}

	orderByStr := stmt.OrderByStr
	if orderByStr.Len() > 0 {
		orderByStr.WriteString(", ")
	}

	bufferJoin(orderByStr, colNames, " ASC, ")
	orderByStr.WriteString(" ASC")

	return stmt
}

// Limit generate LIMIT offset, limit statement
func (stmt *Stmt) Limit(limit int, offset ...int) *Stmt {
	stmt.LimitN = limit
	if len(offset) > 0 {
		stmt.Offset = offset[0]
	}
	return stmt
}
