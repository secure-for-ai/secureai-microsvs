package db

import "fmt"

type condOr []Cond

var _ Cond = condOr{}

// Or sets OR conditions
func Or(conds ...Cond) Cond {
	// return condEmpty if no cond is passed in
	length := len(conds)
	if length == 0 {
		return condEmpty{}
	} else if length == 1 {
		return conds[0]
	}
	// make the condition
	var result = make(condOr, 0, length)
	for _, cond := range conds {
		if cond == nil || !cond.IsValid() {
			continue
		}
		result = append(result, cond)
	}

	switch len(result) {
	case 0:
		// return condEmpty if no cond is valid
		return condEmpty{}
	case 1:
		// return result[0] if only one cond is valid
		return result[0]
	default:
		return result
	}
}

// WriteTo implments Cond
func (or condOr) WriteTo(w Writer) error {
	length := len(or) - 1
	for i, cond := range or {
		var wrap bool
		switch cond.(type) {
		case condAnd, expr:
			wrap = true
			//case Eq:
			//	wrap = (len(cond.(Eq)) > 1)
			//case Neq:
			//	wrap = (len(cond.(Neq)) > 1)
		}

		if wrap {
			fmt.Fprint(w, "(")
		}

		err := cond.WriteTo(w)
		if err != nil {
			return err
		}

		if wrap {
			fmt.Fprint(w, ")")
		}

		if i != length {
			fmt.Fprint(w, " OR ")
		}
	}

	return nil
}

func (or condOr) And(conds ...Cond) Cond {
	return And(append([]Cond{or}, conds...)...)
}

func (or condOr) Or(conds ...Cond) Cond {
	return Or(append(or, conds...)...)
}

func (or condOr) IsValid() bool {
	return len(or) > 1
}
