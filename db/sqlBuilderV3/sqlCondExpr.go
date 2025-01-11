package sqlBuilderV3

import (
	"sync"
)

type condExpr struct {
	sql *stringWriter
	args []interface{}
}

var _ Cond = &condExpr{}
var condExprPool = sync.Pool{
	New: func() interface{} {
		return new(condExpr)
	},
}

// Expr generate customize SQL
func Expr(sql string, args ...interface{}) *condExpr {
	var expr = condExprPool.Get().(*condExpr)
	expr.sql = bufPool.Get().(*stringWriter)
	expr.set(sql, args...)
	return expr
}

func CondExpr(sql string, args ...interface{}) Cond {
	if len(sql) == 0 {
		return CondEmpty
	}
	return Expr(sql, args...)
}

func (expr *condExpr) set(sql string, args ...interface{}) {
	expr.sql.WriteString(sql)
	expr.args = append(expr.args[:0], args...)
}

func (expr *condExpr) WriteTo(w *Writer) {
	w.WriteString(expr.String())
	w.Append(expr.args...)
}

// var condPool = sync.Pool{
// 	New: func() interface{} {
// 		res := make([]Cond, 0, 4)
// 		return &res
// 	},
// }

func (expr *condExpr) And(conds ...Cond) (cond Cond) {
	return andOne(expr, conds...)
}

func (expr *condExpr) Or(conds ...Cond) Cond {
	return orOne(expr, conds...)
}

func (expr *condExpr) IsValid() bool {
	return expr.sql.Len() > 0
}

func (expr *condExpr) Reset() {
	expr.sql.Reset()
	expr.args = expr.args[:0]
}

func (expr *condExpr) Destroy() {
	expr.sql.Destroy()
	expr.args = expr.args[:0]
	condExprPool.Put(expr)
}

func (expr *condExpr) String() string {
	return expr.sql.String()
}

func (expr *condExpr) Append(sql string) {
	expr.sql.WriteString(sql)
}

func Eq(sql string, arg interface{}) Cond {
	cond := Expr(sql, arg)
	cond.Append(" = ?")
	return cond
}
