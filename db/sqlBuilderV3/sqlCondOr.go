package sqlBuilderV3

import "sync"

type condOr []Cond

var _ Cond = &condOr{}
var condOrPool = sync.Pool{
	New: func() interface{} {
		conds := make(condOr, 0, 4)
		return &conds
	},
}

// Or generates OR conditions
func Or(conds ...Cond) Cond {
	// return condEmpty if no cond is passed in
	length := len(conds)
	if length == 0 {
		return CondEmpty
	} else if length == 1 {
		if !conds[0].IsValid() {
			return CondEmpty
		}
		return conds[0]
	}

	// make the condition
	var result = condOrPool.Get().(*condOr)
	return orInternal(result, conds...)
}

func OrOne(cond Cond, conds ...Cond) Cond {
	// return condEmpty if no cond is passed in
	length := len(conds)
	if length == 0 {
		return cond
	}

	// make the condition
	var result = condOrPool.Get().(*condOr)
	if cond.IsValid() {
		*result = append(*result, cond)
	}

	return orInternal(result, conds...)
}

func orInternal(result *condOr, conds ...Cond) Cond {
	// add all valid conditions
	for _, cond := range conds {
		if cond == nil || !cond.IsValid() {
			continue
		}
		*result = append(*result, cond)
	}

	switch len(*result) {
	case 0:
		// return condEmpty if no cond is valid.
		// result need to be returned to the Pool.
		condOrPool.Put(result)
		return CondEmpty
	case 1:
		// return result[0] if only one cond is valid.
		// result need to be returned to the Pool,
		// and its length resets to 0.
		retVal := (*result)[0]
		*result = (*result)[:0]
		condOrPool.Put(result)
		return retVal
	default:
		return result
	}
}

// WriteTo implments Cond
func (or *condOr) WriteTo(w *Writer) {
	length := len(*or) - 1
	for i, cond := range *or {
		var wrap bool
		switch cond.(type) {
		case *condAnd, *condExpr:
			wrap = true
			//case Eq:
			//	wrap = (len(cond.(Eq)) > 1)
			//case Neq:
			//	wrap = (len(cond.(Neq)) > 1)
		}

		if wrap {
			w.WriteByte('(')
		}

		cond.WriteTo(w)

		if wrap {
			w.WriteByte(')')
		}

		if i != length {
			w.WriteString(" OR ")
		}
	}
}

func (or *condOr) And(conds ...Cond) Cond {
	return AndOne(or, conds...)
}

func (or *condOr) Or(conds ...Cond) Cond {
	return OrOne(or, conds...)
}

func (or *condOr) IsValid() bool {
	return len(*or) > 1
}

func (or *condOr) Reset() {
	if len(*or) > 0 {
		// we don't destroy cond recursively, as underlying cond can be
		// used by other sql as well. No need to worry about the dirty data
		// in the array as it is just a reference and will be overwritten
		// next time.
		*or = (*or)[:0]
	}
}

func (or *condOr) Destroy() {
	or.Reset()
	condOrPool.Put(or)
}
