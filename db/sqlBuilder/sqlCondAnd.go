package sqlBuilder

import "fmt"

type condAnd []Cond

var _ Cond = condAnd{}

// And generates AND conditions
func And(conds ...Cond) Cond {
	// return condEmpty if no cond is passed in
	length := len(conds)
	if length == 0 {
		return condEmpty{}
	} else if length == 1 {
		return conds[0]
	}
	// make the condition
	var result = make(condAnd, 0, length)
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

func (and condAnd) WriteTo(w Writer) error {
	length := len(and) - 1
	for i, cond := range and {
		var wrap bool
		switch cond.(type) {
		case condOr, expr:
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
			fmt.Fprint(w, " AND ")
		}
	}

	return nil
}

func (and condAnd) And(conds ...Cond) Cond {
	return And(append(and, conds...)...)
}

func (and condAnd) Or(conds ...Cond) Cond {
	return Or(append([]Cond{and}, conds...)...)
}

func (and condAnd) IsValid() bool {
	return len(and) > 1
}
