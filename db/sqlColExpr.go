package db

import "fmt"

type colParams struct {
	ColNames []string
	Args     []interface{}
}

func (exprs *colParams) addParam(colName string, arg interface{}) {
	exprs.ColNames = append(exprs.ColNames, colName)
	exprs.Args = append(exprs.Args, arg)
}

// int must greater than or equal to 0
func (exprs *colParams) extend(size int) {
	nLen := len(exprs.ColNames)
	nCap := cap(exprs.ColNames)
	newLen := nLen + size

	if newLen > nCap {
		newColNames := make([]string, nCap, newLen)
		copy(newColNames, exprs.ColNames)
		newArgs := make([]interface{}, nCap, newLen)
		copy(newArgs, exprs.Args)
		exprs.ColNames = newColNames
		exprs.Args = newArgs
	}
}

func (exprs *colParams) setColNames(cols []string) {
	exprs.ColNames = make([]string, len(cols))
	copy(exprs.ColNames, cols)
}

func (exprs *colParams) setArgs(args []interface{}) {
	exprs.Args = make([]interface{}, len(args))
	copy(exprs.Args, args)
}

func (exprs *colParams) writeNameArgs(w Writer) error {
	for i, colName := range exprs.ColNames {
		if _, err := fmt.Fprint(w, colName); err != nil {
			return err
		}
		if _, err := fmt.Fprint(w, " = "); err != nil {
			return err
		}

		switch arg := exprs.Args[i].(type) {
		case *SQLStmt:
			if _, err := fmt.Fprint(w, "("); err != nil {
				return err
			}
			if err := arg.WriteTo(w); err != nil {
				return err
			}
			if _, err := fmt.Fprint(w, "("); err != nil {
				return err
			}
		case expr:
			if err := arg.WriteTo(w); err != nil {
				return err
			}
		default:
			w.Append(exprs.Args[i])
		}

		if i+1 != len(exprs.ColNames) {
			if _, err := fmt.Fprint(w, ","); err != nil {
				return err
			}
		}
	}
	return nil
}
