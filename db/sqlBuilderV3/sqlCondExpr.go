package sqlBuilderV3

import (
	"sync"

	"github.com/secure-for-ai/secureai-microsvs/db"
)

type condExpr struct {
	sql  *stringWriter
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

func ExprInc(col string, args ...interface{}) *condExpr {
	var para interface{} = 1
	if len(args) > 0 {
		para = args[0]
	}
	cond := Expr(col, para)
	cond.appendSql(" = ")
	cond.appendSql(col)
	cond.appendSql(" + ")
	cond.appendSql(db.Para)

	return cond
}

func ExprDec(col string, args ...interface{}) *condExpr {
	var para interface{} = 1
	if len(args) > 0 {
		para = args[0]
	}
	cond := Expr(col, para)
	cond.appendSql(" = ")
	cond.appendSql(col)
	cond.appendSql(" - ")
	cond.appendSql(db.Para)

	return cond
}

func ExprSet(col string, val string, args ...interface{}) *condExpr {
	cond := Expr(col, args...)
	cond.appendSql(" = ")
	cond.appendSql(val)
	return cond
}

// set col append
func (expr *condExpr) appendSql(sql string) {
	expr.sql.WriteString(sql)
}

func (expr *condExpr) appendArgs(args ...interface{}) {
	expr.args = append(expr.args, args...)
}

func (expr *condExpr) appendArg(arg interface{}) {
	expr.args = append(expr.args, arg)
}

func ExprEq(sql string, arg interface{}) *condExpr {
	cond := Expr(sql, arg)
	cond.appendSql(" = ")
	cond.appendSql(db.Para)

	return cond
}

func (setCols *condExpr) appendExpr(e *condExpr) {
	if setCols.IsValid() {
		setCols.appendSql(",")
	}
	setCols.appendSql(e.String())
	setCols.appendArgs(e.args...)
}

func (setCols *condExpr) appendEq(col string, arg interface{}) {
	if setCols.IsValid() {
		setCols.appendSql(",")
	}
	setCols.appendSql(col)
	setCols.appendSql(" = ")
	setCols.appendSql(db.Para)
	setCols.appendArg(arg)
}

func (setCols *condExpr) appendInc(col string, args ...interface{}) {
	var para interface{} = 1
	if len(args) > 0 {
		para = args[0]
	}

	if setCols.IsValid() {
		setCols.appendSql(",")
	}
	setCols.appendSql(col)
	setCols.appendSql(" = ")
	setCols.appendSql(col)
	setCols.appendSql(" + ")
	setCols.appendSql(db.Para)
	setCols.appendArg(para)
}

func (setCols *condExpr) appendDec(col string, args ...interface{}) {
	var para interface{} = 1
	if len(args) > 0 {
		para = args[0]
	}

	if setCols.IsValid() {
		setCols.appendSql(",")
	}
	setCols.appendSql(col)
	setCols.appendSql(" = ")
	setCols.appendSql(col)
	setCols.appendSql(" - ")
	setCols.appendSql(db.Para)
	setCols.appendArg(para)
}

func (setCols *condExpr) appendSet(col string, val string, args ...interface{}) {
	if setCols.IsValid() {
		setCols.appendSql(",")
	}
	setCols.appendSql(col)
	setCols.appendSql(" = ")
	setCols.appendSql(val)
	setCols.appendArgs(args...)
}
