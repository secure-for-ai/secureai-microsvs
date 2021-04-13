package db

import (
	"fmt"
	"github.com/secure-for-ai/secureai-microsvs/util"
	"reflect"
	"sort"
	"strings"
)

type SQLType int
type SQLSchema int

const (
	SQLNull SQLType = iota
	SQLInsert
	SQLDelete
	SQLUpdate
	SQLSelect
	SQLUpsert
	SQLTag                = "db"
	SQLPara               = "??"
	SQLPOSTGRES SQLSchema = iota
	SQLMYSQL
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

type SQLTable interface {
	GetTableName() string
}

type fromItem interface {
	itemName() string
	aliasName() string
	setAliasName(string)
	writeTo(writer Writer) error
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

func (from *fromTable) writeTo(w Writer) error {
	fmt.Fprint(w, from.tableName)
	if len(from.alias) > 0 {
		fmt.Fprint(w, " AS ", from.alias)
	}
	return nil
}

type fromStmt struct {
	stmt  *SQLStmt
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

func (from *fromStmt) writeTo(w Writer) error {
	fmt.Fprint(w, "(")
	if err := from.stmt.WriteTo(w); err != nil {
		return err
	}
	fmt.Fprint(w, "(")
	if len(from.alias) > 0 {
		fmt.Fprint(w, ") AS ", from.alias)
	} else {
		fmt.Fprint(w, ")")
	}
	return nil
}

type SQLStmt struct {
	RefTable *SQLTable

	tableInto string
	tableFrom []fromItem

	cond Cond

	GroupByStr string
	having     Cond
	OrderByStr string

	Offset int
	LimitN int

	InsertCols   []string
	InsertValues [][]interface{}
	//isInsertBulk bool

	SetCols colParams

	SelectCols []string

	sqlType      SQLType
	insertSelect *SQLStmt
}

func Insert(data ...interface{}) *SQLStmt {
	return SQL().Insert(data...)
}

func InsertBulk(data interface{}) *SQLStmt {
	return SQL().InsertBulk(data)
}

func Delete(data ...interface{}) *SQLStmt {
	return SQL().Delete(data...)
}

func Update(data ...interface{}) *SQLStmt {
	return SQL().Update(data...)
}

func Select(data ...interface{}) *SQLStmt {
	return SQL().Select(data...)
}

func SQL() *SQLStmt {
	stmt := &SQLStmt{}
	stmt.Init()
	return stmt
}

// Init reset all the statement's fields
func (stmt *SQLStmt) Init() {
	stmt.RefTable = nil

	stmt.tableInto = ""
	stmt.tableFrom = make([]fromItem, 0, 2)
	stmt.cond = condEmpty{}

	stmt.GroupByStr = ""
	stmt.having = condEmpty{}
	stmt.OrderByStr = ""

	stmt.Offset = 0
	stmt.LimitN = 0

	stmt.InsertCols = []string{}
	stmt.InsertValues = [][]interface{}{}
	//stmt.isInsertBulk = false

	stmt.SetCols = colParams{}

	stmt.SelectCols = []string{}

	stmt.sqlType = SQLNull
	stmt.insertSelect = nil
}

// TableName returns the table name
//func (stmt *SQLStmt) TableName() string {
//	if stmt.sqlType == SQLInsert {
//		return stmt.tableIntro
//	}
//	return stmt.tableFrom[0].itemName()
//}

func buildColumns(column interface{}) []string {
	v := util.ReflectValue(column)
	vType := v.Type()
	if vType.Kind() == reflect.Struct {
		numField := v.NumField()
		// avoid extend the slice cap which causes memory reallocation
		colNames := make([]string, numField)
		for i, il := 0, v.NumField(); i < il; i++ {
			// Get column name, tag start with "pg" or the field Name
			var colName string
			fieldInfo := vType.Field(i)
			if colName = fieldInfo.Tag.Get(SQLTag); colName == "" {
				colName = vType.Field(i).Name
			}
			colNames[i] = colName
		}
		return colNames
	}
	return nil
}

func buildValues(curData interface{}) []interface{} {
	v := util.ReflectValue(curData)
	vType := v.Type()
	if vType.Kind() == reflect.Struct {

		numField := v.NumField()
		values := make([]interface{}, numField)
		for i, il := 0, v.NumField(); i < il; i++ {
			// Get value
			var val interface{}
			fieldValue := v.Field(i)
			val = fieldValue.Interface()
			values[i] = Expr(SQLPara, val)
		}
		return values
	}
	return nil
}

// Into sets insert table name
func (stmt *SQLStmt) IntoTable(table interface{}) *SQLStmt {
	switch table.(type) {
	case SQLTable:
		stmt.tableInto = table.(SQLTable).GetTableName()
	case string:
		stmt.tableInto = table.(string)
	}
	return stmt
}

func (stmt *SQLStmt) IntoColumns(column interface{}, cols ...string) *SQLStmt {
	switch column.(type) {
	case []string:
		stmt.InsertCols = append(stmt.InsertCols, column.([]string)...)
	case Columns:
		stmt.InsertCols = append(stmt.InsertCols, column.(Columns)...)
	case string:
		stmt.InsertCols = append(stmt.InsertCols, column.(string))
		stmt.InsertCols = append(stmt.InsertCols, cols...)
	default:
		InsertCols := buildColumns(column)
		if InsertCols != nil {
			stmt.InsertCols = append(stmt.InsertCols, InsertCols...)
		}
	}
	return stmt
}

// Values store the insertion data, optimized for one record and support bulk insertion as well
func (stmt *SQLStmt) Values(data ...interface{}) *SQLStmt {
	switch len(data) {
	case 0:
		return stmt
	case 1:
		//if len(stmt.InsertValues) >= 1 {
		//	stmt.isInsertBulk = true
		//}
		curData := data[0]
		switch curData.(type) {
		case []interface{}:
			dataArray := curData.([]interface{})
			values := make([]interface{}, 0, len(dataArray))
			for _, el := range dataArray {
				if e, ok := el.(expr); ok {
					values = append(values, e)
				} else {
					values = append(values, Expr(SQLPara, e))
				}
			}
			stmt.InsertValues = append(stmt.InsertValues, values)
		case Map:
			dataMap := curData.(Map)
			InsertCols := make([]string, 0, len(dataMap))
			InsertValues := make([]interface{}, 0, len(dataMap))
			for _, col := range dataMap.sortedKeys() {
				InsertCols = append(InsertCols, col)
				val := dataMap[col]
				if e, ok := val.(expr); ok {
					InsertValues = append(InsertValues, e)
				} else {
					InsertValues = append(InsertValues, Expr(SQLPara, val))
				}
			}
			stmt.InsertCols = InsertCols
			stmt.InsertValues = append(stmt.InsertValues, InsertValues)
		case *SQLStmt:
			stmt.insertSelect = curData.(*SQLStmt)
		default:
			if len(stmt.InsertCols) == 0 {
				if columns := buildColumns(curData); columns != nil {
					stmt.InsertCols = columns
				}
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

func (stmt *SQLStmt) ValuesBulk(data interface{}) *SQLStmt {
	dataR := util.ReflectValue(data)
	dataType := dataR.Kind()
	if dataType != reflect.Slice && dataType != reflect.Array {
		return stmt
	}
	return stmt.valuesBulkInternal(&dataR)
}

func (stmt *SQLStmt) valuesBulkInternal(data *reflect.Value) *SQLStmt {
	dataLen := data.Len()
	if dataLen == 0 {
		return stmt
	}

	//if dataLen > 1 || len(stmt.InsertValues) >= 1 {
	//	stmt.isInsertBulk = true
	//}

	// update the insert columns
	data0 := data.Index(0).Interface()
	switch data0.(type) {
	case Map:
		dataMap := data0.(Map)
		InsertCols := make([]string, 0, len(dataMap))
		for _, col := range dataMap.sortedKeys() {
			InsertCols = append(InsertCols, col)
		}
		stmt.InsertCols = InsertCols
	default:
		InsertCols := buildColumns(data0)
		if InsertCols != nil {
			stmt.InsertCols = append(stmt.InsertCols, InsertCols...)
		}
	}

	// loading the data
	InsertValues := make([][]interface{}, 0, dataLen)
	for i := 0; i < dataLen; i++ {
		curData := data.Index(i).Interface()
		switch curData.(type) {
		case []interface{}:
			dataArray := curData.([]interface{})
			values := make([]interface{}, 0, len(dataArray))
			for _, el := range dataArray {
				if e, ok := el.(expr); ok {
					values = append(values, e)
				} else {
					values = append(values, Expr(SQLPara, e))
				}
			}
			InsertValues = append(InsertValues, values)
		case Map:
			dataMap := curData.(Map)
			values := make([]interface{}, 0, len(dataMap))
			for _, col := range stmt.InsertCols {
				val := dataMap[col]
				if e, ok := val.(expr); ok {
					values = append(values, e)
				} else {
					values = append(values, Expr(SQLPara, val))
				}
			}
			InsertValues = append(InsertValues, values)
		case *SQLStmt:
			stmt.insertSelect = curData.(*SQLStmt)
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

func (stmt *SQLStmt) SelectColumns(column interface{}, cols ...string) *SQLStmt {
	switch column.(type) {
	case []string:
		stmt.SelectCols = append(stmt.SelectCols, column.([]string)...)
	case Columns:
		stmt.SelectCols = append(stmt.SelectCols, column.(Columns)...)
	case string:
		stmt.SelectCols = append(stmt.SelectCols, column.(string))
		stmt.SelectCols = append(stmt.SelectCols, cols...)
	default:
		SelectCols := buildColumns(column)
		if SelectCols != nil {
			stmt.SelectCols = append(stmt.SelectCols, SelectCols...)
		}
	}
	return stmt
}

// From sets from subject(can be a table name in string or a builder pointer) and its alias
func (stmt *SQLStmt) From(subject interface{}, alias ...string) *SQLStmt {
	var from fromItem
	switch subject.(type) {
	case *SQLStmt:
		//subquery should be a select statement
		from = &fromStmt{
			subject.(*SQLStmt),
			"",
		}
	case SQLTable:
		from = &fromTable{
			subject.(SQLTable).GetTableName(),
			"",
		}
	case string:
		from = &fromTable{
			subject.(string),
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
func (stmt *SQLStmt) Insert(data ...interface{}) *SQLStmt {
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
	if stmt.sqlType == SQLNull {
		stmt.sqlType = SQLInsert
	}
	return stmt
}

// Insert SQL
func (stmt *SQLStmt) InsertBulk(data interface{}) *SQLStmt {
	if stmt.sqlType == SQLNull {
		stmt.sqlType = SQLInsert
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
func (stmt *SQLStmt) Delete(data ...interface{}) *SQLStmt {
	l := len(data)
	if l >= 1 {
		stmt.From(data[0])
	}
	if l >= 2 {
		stmt.And(data[1], data[2:]...)
	}
	if stmt.sqlType == SQLNull {
		stmt.sqlType = SQLDelete
	}
	return stmt
}

// Update
func (stmt *SQLStmt) Update(data ...interface{}) *SQLStmt {
	l := len(data)
	if l >= 1 {
		stmt.Set(data[0])
		stmt.From(data[0])
	}
	if l >= 2 {
		stmt.And(data[1], data[2:]...)
	}
	if stmt.sqlType == SQLNull {
		stmt.sqlType = SQLUpdate
	}
	return stmt
}

// Select SQL
func (stmt *SQLStmt) Select(data ...interface{}) *SQLStmt {
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
	if stmt.sqlType == SQLNull {
		stmt.sqlType = SQLSelect
	}
	return stmt
}

// Incr Generate  "Update ... Set column = column + arg" statement
func (stmt *SQLStmt) Incr(col string, arg ...interface{}) *SQLStmt {
	var para interface{} = 1
	if len(arg) > 0 {
		para = arg[0]
	}
	stmt.SetCols.addParam(col, Expr(col+" + "+SQLPara, para))
	return stmt
}

// Decr Generate  "Update ... Set column = column - arg" statement
func (stmt *SQLStmt) Decr(col string, arg ...interface{}) *SQLStmt {
	var para interface{} = 1
	if len(arg) > 0 {
		para = arg[0]
	}
	stmt.SetCols.addParam(col, Expr(col+" - "+SQLPara, para))
	return stmt
}

// setExpr Generate  "Update ... Set column = {expr}" statement
// if you want to use writeTo internal builtin functions without parameters like NOW(),
// then you'd better to call Set(col, Expr("Now()"))
// Todo support expr as SQLStmt
func (stmt *SQLStmt) setExpr(col string, expr interface{}, args ...interface{}) *SQLStmt {
	if e, ok := expr.(string); ok {
		if len(args) > 0 {
			// set("col", "col||??", "test") => writeTo: col = col||??, args: "test"
			stmt.SetCols.addParam(col, Expr(e, args...))
		} else {
			// set("col", "test") => writeTo: col = ??, args: "test"
			// equivalent to set("col", SQLPara, "test")
			stmt.SetCols.addParam(col, Expr(SQLPara, e))
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
func (stmt *SQLStmt) setMap(exprs Map) *SQLStmt {
	// avoid extend the slice cap which causes memory reallocation
	stmt.SetCols.extend(len(exprs))
	for col, val := range exprs {
		if e, ok := val.(expr); ok {
			stmt.SetCols.addParam(col, e)
		} else {
			stmt.SetCols.addParam(col, Expr(SQLPara, val))
		}
	}
	return stmt
}

func (stmt *SQLStmt) setStruct(data interface{}) *SQLStmt {
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
			if colName = fieldInfo.Tag.Get(SQLTag); colName == "" {
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

			stmt.SetCols.addParam(colName, Expr(SQLPara, val))
		}
	}
	return stmt
}

func (stmt *SQLStmt) Set(data interface{}, args ...interface{}) *SQLStmt {
	switch data.(type) {
	case string:
		argLen := len(args)
		if argLen >= 1 {
			stmt.setExpr(data.(string), args[0], args[1:]...)
		} else {
			// Todo Raise Error
		}
	case Map:
		stmt.setMap(data.(Map))
	default:
		// assume the input is either a struct ptr or a struct
		stmt.setStruct(data)
	}
	return stmt
}

func (stmt *SQLStmt) Where(query interface{}, args ...interface{}) *SQLStmt {
	return stmt.catCond(&stmt.cond, And, query, args...)
}

// concat an existing Cond and a new Cond statement with Op
func (stmt *SQLStmt) catCond(c *Cond, OpFunc func(cond ...Cond) Cond, query interface{}, args ...interface{}) *SQLStmt {
	switch query.(type) {
	case string:
		cond := Expr(query.(string), args...)
		*c = OpFunc(*c, cond)
	case Map:
		queryMap := query.(Map)
		conds := make([]Cond, 0, len(queryMap)+1)
		conds = append(conds, *c)
		for _, k := range queryMap.sortedKeys() {
			conds = append(conds, Expr(k+" = "+SQLPara, queryMap[k]))
		}
		*c = OpFunc(conds...)
	case Cond:
		conds := make([]Cond, 0, len(args)+2)
		conds = append(conds, *c)
		conds = append(conds, query.(Cond))
		for _, v := range args {
			if vv, ok := v.(Cond); ok {
				conds = append(conds, vv)
			}
		}
		*c = OpFunc(conds...)
	default:
		// TODO: not support condition type
	}
	return stmt
}

// And add Where & and statement
func (stmt *SQLStmt) And(query interface{}, args ...interface{}) *SQLStmt {
	return stmt.catCond(&stmt.cond, And, query, args...)
}

// Or add Where & Or statement
func (stmt *SQLStmt) Or(query interface{}, args ...interface{}) *SQLStmt {
	return stmt.catCond(&stmt.cond, Or, query, args...)
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
func (stmt *SQLStmt) GroupBy(keys ...string) *SQLStmt {
	if len(keys) == 0 {
		return stmt
	}
	if len(stmt.GroupByStr) > 0 {
		stmt.GroupByStr += ", "
	}
	stmt.GroupByStr = strings.Join(keys, ", ")
	return stmt
}

// GroupBy generate "Having conditions" statement
func (stmt *SQLStmt) Having(query interface{}, args ...interface{}) *SQLStmt {
	return stmt.catCond(&stmt.having, And, query, args...)
}

// GroupBy generate "Having conditions" statement && conditions
func (stmt *SQLStmt) HavingAnd(query interface{}, args ...interface{}) *SQLStmt {
	return stmt.catCond(&stmt.having, And, query, args...)
}

// GroupBy generate "Having conditions" statement || conditions
func (stmt *SQLStmt) HavingOr(query interface{}, args ...interface{}) *SQLStmt {
	return stmt.catCond(&stmt.having, Or, query, args...)
}

// OrderBy generate "Order By order" statement
func (stmt *SQLStmt) OrderBy(order ...string) *SQLStmt {
	if len(order) == 0 {
		return stmt
	}
	if len(stmt.OrderByStr) > 0 {
		stmt.OrderByStr += ", "
	}

	stmt.OrderByStr += strings.Join(order, ", ") // statement.ReplaceQuote(order) pq.QuoteIdentifier()
	return stmt
}

// Desc generate `ORDER BY xx DESC`
func (stmt *SQLStmt) Desc(colNames ...string) *SQLStmt {
	if len(colNames) == 0 {
		return stmt
	}
	var buf strings.Builder
	if len(stmt.OrderByStr) > 0 {
		fmt.Fprint(&buf, stmt.OrderByStr, ", ")
	}
	fmt.Fprintf(&buf, "%v DESC", strings.Join(colNames, " DESC, "))
	stmt.OrderByStr = buf.String()
	return stmt
}

// Asc generate `ORDER BY xx ASC`
func (stmt *SQLStmt) Asc(colNames ...string) *SQLStmt {
	if len(colNames) == 0 {
		return stmt
	}
	var buf strings.Builder
	if len(stmt.OrderByStr) > 0 {
		fmt.Fprint(&buf, stmt.OrderByStr, ", ")
	}
	fmt.Fprintf(&buf, "%v ASC", strings.Join(colNames, " ASC, "))
	stmt.OrderByStr = buf.String()
	return stmt
}

// Limit generate LIMIT offset, limit statement
func (stmt *SQLStmt) Limit(limit int, offset ...int) *SQLStmt {
	stmt.LimitN = limit
	if len(offset) > 0 {
		stmt.Offset = offset[0]
	}
	return stmt
}

func (stmt *SQLStmt) SQL() (string, []interface{}) {
	return "", []interface{}{}
}
