package sqlBuilderV3

import "sync"

type condExpr struct {
	sql  string
	args []interface{}
}

var _ Cond = &condExpr{}
var condExprPool = sync.Pool{
	New: func() interface{} {
		return new(condExpr)
	},
}

// Expr generate customize SQL
func Expr(sql string, args ...interface{}) Cond {
	if len(sql) == 0 {
		return CondEmpty
	}
	var expr = condExprPool.Get().(*condExpr)
	expr.Set(sql, args...)
	return expr
}

func (expr *condExpr) Set(sql string, args ...interface{}) {
	expr.sql = sql
	// expr.args = args
	// expr.args = expr.args[:0]
	expr.args = append(expr.args[:0], args...)
}

func (expr *condExpr) WriteTo(w *Writer) {
	w.WriteString(expr.sql)
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
	return len(expr.sql) > 0
}

func (expr *condExpr) Reset() {
	expr.sql = ""
	expr.args = expr.args[:0]
}

func (expr *condExpr) Destroy() {
	expr.Reset()
	condExprPool.Put(expr)
}

func Eq(sql string, arg interface{}) Cond {
	if len(sql) == 0 {
		return condEmpty{}
	}
	return Expr(sql+" = ?", arg)
}
