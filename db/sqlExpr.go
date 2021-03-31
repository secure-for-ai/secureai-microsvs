package db

import "fmt"

type expr struct {
	sql  string
	args []interface{}
}

var _ Cond = expr{}

// Expr generate customerize SQL
func Expr(sql string, args ...interface{}) Cond {
	if len(sql) == 0 {
		return condEmpty{}
	}
	return expr{sql, args}
}

//func (expr expr) OpWriteTo(op string, w Writer) error {
//	return expr.WriteTo(w)
//}

func (expr expr) WriteTo(w Writer) error {
	if _, err := fmt.Fprint(w, expr.sql); err != nil {
		return err
	}
	w.Append(expr.args...)
	return nil
}

func (expr expr) And(conds ...Cond) Cond {
	return And(append([]Cond{expr}, conds...)...)
}

func (expr expr) Or(conds ...Cond) Cond {
	return Or(append([]Cond{expr}, conds...)...)
}

func (expr expr) IsValid() bool {
	return len(expr.sql) > 0
}
