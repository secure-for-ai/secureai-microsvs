package sqlBuilderV3

import (
	"bytes"
	"sort"
	"strings"
	"sync"

	// "github.com/goccy/go-reflect"
	"reflect"

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
type Map map[string]interface{}

func (m Map) sortedKeys() []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

type Table interface {
	GetTableName() string
}

type fromItem interface {
	itemName() string
	aliasName() string
	setAliasName(string)
	writeTo(*Writer)
}

type fromTable struct {
	tableName string
	alias     string
}

func (from *fromTable) itemName() string {
	return from.tableName
}

func (from *fromTable) aliasName() string {
	return from.alias
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

type fromStmt struct {
	stmt  *Stmt
	alias string
}

func (from *fromStmt) itemName() string {
	return from.alias
}

func (from *fromStmt) aliasName() string {
	return from.alias
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

type Stmt struct {
	RefTable *Table

	tableInto string
	tableFrom []fromItem

	where Cond

	GroupByStr *bytes.Buffer
	having     Cond
	OrderByStr *bytes.Buffer

	Offset int
	LimitN int

	InsertCols   []string
	InsertValues [][]condExpr
	//isInsertBulk bool

	SetCols colParams

	SelectCols []string

	sqlType      Type
	insertSelect *Stmt
}

func Insert(data ...interface{}) *Stmt {
	return SQL().Insert(data...)
}

func InsertBulk(data interface{}) *Stmt {
	return SQL().InsertBulk(data)
}

func Delete(data ...interface{}) *Stmt {
	return SQL().Delete(data...)
}

func Update(data ...interface{}) *Stmt {
	return SQL().Update(data...)
}

func Select(data ...interface{}) *Stmt {
	return SQL().Select(data...)
}

var stmtPool = sync.Pool{
	New: func() interface{} {
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
	stmt.GroupByStr = new(bytes.Buffer)
	stmt.having = condEmpty{}
	stmt.OrderByStr = new(bytes.Buffer)

	stmt.Offset = 0
	stmt.LimitN = 0

	stmt.InsertCols = []string{}
	stmt.InsertValues = [][]condExpr{}
	//stmt.isInsertBulk = false
	stmt.SetCols = colParams{}
	stmt.SelectCols = []string{}

	stmt.sqlType = RawType
	stmt.insertSelect = nil
}

func (stmt *Stmt) Reset() {
	stmt.RefTable = nil

	stmt.tableInto = ""
	stmt.tableFrom = stmt.tableFrom[:0]

	// stmt.where.Destroy()
	stmt.where = condEmpty{}
	stmt.GroupByStr.Reset()
	stmt.having = condEmpty{}
	stmt.OrderByStr.Reset()

	stmt.Offset = 0
	stmt.LimitN = 0

	stmt.InsertCols = stmt.InsertCols[:0]
	stmt.InsertValues = stmt.InsertValues[:0]
	stmt.SetCols = colParams{}
	stmt.SelectCols = stmt.SelectCols[:0]

	stmt.sqlType = RawType
	stmt.insertSelect = nil
}

func (stmt *Stmt) Destroy() {
	stmt.Reset()
	stmtPool.Put(stmt)
}

// TableName returns the table name
//func (stmt *SQLStmt) TableName() string {
//	if stmt.sqlType == InsertType {
//		return stmt.tableIntro
//	}
//	return stmt.tableFrom[0].itemName()
//}

// Todo replace with the fast version of sync map
var structColumnCache = sync.Map{}

type stringWriter struct {
	bytes.Buffer
}

func (w *stringWriter) Reset() {
	w.Buffer.Reset()
	bufPool.Put(w)
}

func (w *stringWriter) String() string {
	return util.FastBytesToString(w.Bytes())
}

var bufPool = sync.Pool{
	New: func() interface{} {
		return &stringWriter{}
	},
}

func buildColumns(colNames *[]string, column interface{}) {
	v := util.ReflectValue(column)
	// println(v.String(), v.Type().PkgPath())
	vType := v.Type()
	// println(vType.Name(), vType.PkgPath(), vType.String())
	if vType.Kind() == reflect.Struct {
		// construct the unique name of the struct with zero-copy method
		w := bufPool.Get().(*stringWriter)
		w.WriteString(vType.PkgPath())
		w.WriteByte('.')
		w.WriteString(vType.Name())
		structFullName := w.String()
		structColNames, ok := structColumnCache.Load(structFullName)

		// write the column name of the struct from the cache
		if ok {
			*colNames = append(*colNames, structColNames.([]string)...)
			// free *byte.Buffer
			w.Reset()
			return
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
		*colNames = append(*colNames, tmpColNames...)
		structColumnCache.Store(strings.Clone(structFullName), tmpColNames)
		// free *byte.Buffer
		w.Reset()
	}
}

func buildValues(curData interface{}) []condExpr {
	v := util.ReflectValue(curData)
	vType := v.Type()
	if vType.Kind() == reflect.Struct {

		numField := v.NumField()
		values := make([]condExpr, numField)
		for i, il := 0, numField; i < il; i++ {
			// Get value
			var val interface{}
			fieldValue := v.Field(i)
			val = fieldValue.Interface()
			values[i].Set(db.Para, val)
		}
		return values
	}
	return nil
}

// Into sets insert table name
func (stmt *Stmt) IntoTable(table interface{}) *Stmt {
	switch table := table.(type) {
	case Table:
		stmt.tableInto = table.GetTableName()
	case string:
		stmt.tableInto = table
	}
	return stmt
}

func (stmt *Stmt) IntoColumns(column interface{}, cols ...string) *Stmt {
	switch column := column.(type) {
	case []string:
		stmt.InsertCols = append(stmt.InsertCols, column...)
	case Columns:
		stmt.InsertCols = append(stmt.InsertCols, column...)
	case string:
		stmt.InsertCols = append(stmt.InsertCols, column)
		stmt.InsertCols = append(stmt.InsertCols, cols...)
	default:
		buildColumns(&stmt.InsertCols, column)
		// if InsertCols != nil {
		// 	stmt.InsertCols = append(stmt.InsertCols, InsertCols...)
		// }
	}
	return stmt
}

// Values store the insertion data, optimized for one record and support bulk insertion as well
func (stmt *Stmt) Values(data ...interface{}) *Stmt {
	switch len(data) {
	case 0:
		return stmt
	case 1:
		//if len(stmt.InsertValues) >= 1 {
		//	stmt.isInsertBulk = true
		//}
		curData := data[0]
		switch curData := curData.(type) {
		case []interface{}:
			values := make([]condExpr, len(curData))
			for i, el := range curData {
				if e, ok := el.(condExpr); ok {
					values[i] = e
				} else {
					values[i].Set(db.Para, e)
				}
			}
			stmt.InsertValues = append(stmt.InsertValues, values)
		case Map:
			InsertCols := make([]string, 0, len(curData))
			InsertValues := make([]condExpr, len(curData))
			for i, col := range curData.sortedKeys() {
				InsertCols = append(InsertCols, col)
				val := curData[col]
				if e, ok := val.(condExpr); ok {
					InsertValues[i] = e
				} else {
					InsertValues[i].Set(db.Para, val)
				}
			}
			stmt.InsertCols = InsertCols
			stmt.InsertValues = append(stmt.InsertValues, InsertValues)
		case *Stmt:
			stmt.insertSelect = curData
		default:
			if len(stmt.InsertCols) == 0 {
				buildColumns(&stmt.InsertCols, curData)
			}
			values := buildValues(curData)
			if values != nil {
				stmt.InsertValues = append(stmt.InsertValues, values)
			}
		}
	default:
		return stmt.ValuesBulk(data)
	}

	return stmt
}

func (stmt *Stmt) ValuesBulk(data interface{}) *Stmt {
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

	//if dataLen > 1 || len(stmt.InsertValues) >= 1 {
	//	stmt.isInsertBulk = true
	//}

	// update the insert columns
	data0 := data.Index(0).Interface()
	switch data0 := data0.(type) {
	case Map:
		// InsertCols := make([]string, 0, len(data0))
		// InsertCols = append(InsertCols, data0.sortedKeys()...)
		// stmt.InsertCols = InsertCols
		stmt.InsertCols = append(stmt.InsertCols, data0.sortedKeys()...)
	default:
		buildColumns(&stmt.InsertCols, data0)
		// if InsertCols != nil {
		// 	stmt.InsertCols = append(stmt.InsertCols, InsertCols...)
		// }
	}

	// loading the data
	InsertValues := make([][]condExpr, 0, dataLen)
	for i := 0; i < dataLen; i++ {
		curData := data.Index(i).Interface()
		switch curData := curData.(type) {
		case []interface{}:
			values := make([]condExpr, len(curData))
			for j, el := range curData {
				if e, ok := el.(condExpr); ok {
					values[j] = e
				} else {
					values[j].Set(db.Para, e)
				}
			}
			InsertValues = append(InsertValues, values)
		case Map:
			values := make([]condExpr, len(curData))
			for j, col := range stmt.InsertCols {
				val := curData[col]
				if e, ok := val.(condExpr); ok {
					values[j] = e
				} else {
					values[j].Set(db.Para, val)
				}
			}
			InsertValues = append(InsertValues, values)
		case *Stmt:
			stmt.insertSelect = curData
		default:
			values := buildValues(curData)
			if values != nil {
				InsertValues = append(InsertValues, values)
			}
		}
	}
	stmt.InsertValues = append(stmt.InsertValues, InsertValues...)

	return stmt
}

func (stmt *Stmt) SelectColumns(column interface{}, cols ...string) *Stmt {
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
		// if SelectCols != nil {
		// 	stmt.SelectCols = append(stmt.SelectCols, SelectCols...)
		// }
	}
	return stmt
}

// From sets from subject(can be a table name in string or a builder pointer) and its alias
func (stmt *Stmt) From(subject interface{}, alias ...string) *Stmt {
	var from fromItem
	switch subject := subject.(type) {
	case *Stmt:
		//subquery should be a select statement
		from = &fromStmt{
			subject,
			"",
		}
	case Table:
		from = &fromTable{
			subject.GetTableName(),
			"",
		}
	case string:
		from = &fromTable{
			subject,
			"",
		}
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
func (stmt *Stmt) Insert(data ...interface{}) *Stmt {
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

// Insert SQL
func (stmt *Stmt) InsertBulk(data interface{}) *Stmt {
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

	//switch len(data) {
	//case 0:
	//	break
	//default:
	//	// if data is a double array, you need to call IntoColumns afterwards.
	//	// Otherwise the order of the values should be the same order of the columns in the table.
	//	// Support Bulk Insertion
	//	s := reflect.ValueOf(data)
	//	s.Len()
	//	stmt.ValuesBulk(data)
	//	stmt.IntoTable(data[0])
	//}
	return stmt
}

// Delete sets delete SQL
func (stmt *Stmt) Delete(data ...interface{}) *Stmt {
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
func (stmt *Stmt) Update(data ...interface{}) *Stmt {
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
func (stmt *Stmt) Select(data ...interface{}) *Stmt {
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
func (stmt *Stmt) Incr(col string, arg ...interface{}) *Stmt {
	var para interface{} = 1
	if len(arg) > 0 {
		para = arg[0]
	}
	stmt.SetCols.addParam(col, Expr(col+" + "+db.Para, para))
	return stmt
}

// Decr Generate  "Update ... Set column = column - arg" statement
func (stmt *Stmt) Decr(col string, arg ...interface{}) *Stmt {
	var para interface{} = 1
	if len(arg) > 0 {
		para = arg[0]
	}
	stmt.SetCols.addParam(col, Expr(col+" - "+db.Para, para))
	return stmt
}

// setExpr Generate  "Update ... Set column = {expr}" statement
// if you want to use writeTo internal builtin functions without parameters like NOW(),
// then you'd better to call Set(col, Expr("Now()"))
// Todo support expr as SQLStmt
func (stmt *Stmt) setExpr(col string, expr interface{}, args ...interface{}) *Stmt {
	if e, ok := expr.(string); ok {
		if len(args) > 0 {
			// set("col", "col||??", "test") => writeTo: col = col||??, args: "test"
			stmt.SetCols.addParam(col, Expr(e, args...))
		} else {
			// set("col", "test") => writeTo: col = ??, args: "test"
			// equivalent to set("col", Para, "test")
			stmt.SetCols.addParam(col, Expr(db.Para, e))
		}
	} else {
		stmt.SetCols.addParam(col, expr)
	}
	return stmt
}

// setMap Generate  "Update ... Set col1 = {expr1}, col1 = {expr2}" statement
// {"username": "bob", "age": 10, "createTime": Expr("Now()"} =>
// SQL: username = ?? , age = ??, createTime = NOW()
// Args: ["bob", 10]
// Todo support expr as SQLStmt
func (stmt *Stmt) setMap(exprs Map) *Stmt {
	// avoid extend the slice cap which causes memory reallocation
	stmt.SetCols.extend(len(exprs))
	for col, val := range exprs {
		if e, ok := val.(condExpr); ok {
			stmt.SetCols.addParam(col, e)
		} else {
			stmt.SetCols.addParam(col, Expr(db.Para, val))
		}
	}
	return stmt
}

func (stmt *Stmt) setStruct(data interface{}) *Stmt {
	// check whether data is struct
	// reflect the exact value of the data regardless of whether it's a ptr or struct
	v := util.ReflectValue(data)
	vType := v.Type()
	if vType.Kind() == reflect.Struct {

		numField := v.NumField()
		// avoid extend the slice cap which causes memory reallocation
		stmt.SetCols.extend(numField)

		for i, il := 0, v.NumField(); i < il; i++ {
			// Get column name, tag start with "pg" or the field Name
			var colName string
			fieldInfo := vType.Field(i)
			if colName = fieldInfo.Tag.Get(db.Tag); colName == "" {
				colName = vType.Field(i).Name
			}

			// Get value
			var val interface{}
			fieldValue := v.Field(i)
			val = fieldValue.Interface()

			/*fieldType := reflect.TypeOf(fieldValue.Interface())
			switch fieldType.Kind() {
			case reflect.Bool:
				val = fieldValue.Bool()
			case reflect.String:
				val = fieldValue.String()
			case reflect.Int8, reflect.Int16, reflect.Int, reflect.Int32, reflect.Int64:
				val = fieldValue.Int()
			case reflect.Float32, reflect.Float64:
				val = fieldValue.Float()
			case reflect.Uint8, reflect.Uint16, reflect.Uint, reflect.Uint32, reflect.Uint64:
				val = fieldValue.Uint()
			default:
				val = fieldValue.Interface()
			}*/

			stmt.SetCols.addParam(colName, Expr(db.Para, val))
		}
	}
	return stmt
}

func (stmt *Stmt) Set(data interface{}, args ...interface{}) *Stmt {
	switch data := data.(type) {
	case string:
		argLen := len(args)
		if argLen >= 1 {
			stmt.setExpr(data, args[0], args[1:]...)
		} else {
			// Todo Raise Error
		}
	case Map:
		stmt.setMap(data)
	default:
		// assume the input is either a struct ptr or a struct
		stmt.setStruct(data)
	}
	return stmt
}

func (stmt *Stmt) Where(query interface{}, args ...interface{}) *Stmt {
	stmt.where = stmt.catCond(stmt.where, And, query, args...)
	return stmt
}

// concat an existing Cond and a new Cond statement with Op
func (stmt *Stmt) catCond(c Cond, OpFunc func(cond ...Cond) Cond, query interface{}, args ...interface{}) Cond {
	switch query := query.(type) {
	case string:
		cond := Expr(query, args...)
		c = OpFunc(c, cond)
	case Map:
		conds := make([]Cond, 0, len(query)+1)
		conds = append(conds, c)
		for _, k := range query.sortedKeys() {
			conds = append(conds, Expr(k+" = "+db.Para, query[k]))
		}
		c = OpFunc(conds...)
	case Cond:
		conds := make([]Cond, 0, len(args)+2)
		conds = append(conds, c)
		conds = append(conds, query)
		for _, v := range args {
			if vv, ok := v.(Cond); ok {
				conds = append(conds, vv)
			}
		}
		c = OpFunc(conds...)
	default:
		// TODO: not support condition type
	}
	return c
}

// And add Where & and statement
func (stmt *Stmt) And(query interface{}, args ...interface{}) *Stmt {
	stmt.where = stmt.catCond(stmt.where, And, query, args...)
	return stmt
}

// Or add Where & Or statement
func (stmt *Stmt) Or(query interface{}, args ...interface{}) *Stmt {
	stmt.where = stmt.catCond(stmt.where, Or, query, args...)
	return stmt
}

//// In generate "Where column IN (??) " statement
//func (stmt *SQLStmt) In(column string, args ...interface{}) *SQLStmt {
//	in := builder.In(stmt.quote(column), args...)
//	stmt.cond = stmt.cond.And(in)
//	return stmt
//}
//
//// NotIn generate "Where column NOT IN (??) " statement
//func (stmt *SQLStmt) NotIn(column string, args ...interface{}) *SQLStmt {
//	notIn := builder.NotIn(stmt.quote(column), args...)
//	stmt.cond = stmt.cond.And(notIn)
//	return stmt
//}

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
func (stmt *Stmt) Having(query interface{}, args ...interface{}) *Stmt {
	stmt.having = stmt.catCond(stmt.having, And, query, args...)
	return stmt
}

// GroupBy generate "Having conditions" statement && conditions
func (stmt *Stmt) HavingAnd(query interface{}, args ...interface{}) *Stmt {
	stmt.having = stmt.catCond(stmt.having, And, query, args...)
	return stmt
}

// GroupBy generate "Having conditions" statement || conditions
func (stmt *Stmt) HavingOr(query interface{}, args ...interface{}) *Stmt {
	stmt.having = stmt.catCond(stmt.having, Or, query, args...)
	return stmt
}

func bufferJoin(w *bytes.Buffer, elems []string, sep string) {
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
	// stmt.OrderByStr += strings.Join(order, ", ") // statement.ReplaceQuote(order) pq.QuoteIdentifier()
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

func (stmt *Stmt) SQL() (string, []interface{}) {
	return "", []interface{}{}
}
