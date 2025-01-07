package sqlBuilderV3

import (
	// "fmt"
	"sync"
)

type condAnd []Cond

var _ Cond = &condAnd{}
var condAndPool = sync.Pool{
	New: func() interface{} {
		conds := make(condAnd, 0, 4)
		return &conds
	},
}

// And generates AND conditions
func And(conds ...Cond) Cond {
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
	var result = condAndPool.Get().(*condAnd)
	return andInternal(result, conds...)
}

func andOne(cond Cond, conds ...Cond) Cond {
	// return condEmpty if no cond is passed in
	length := len(conds)
	if length == 0 {
		return cond
	}

	// make the condition
	var result = condAndPool.Get().(*condAnd)
	if cond.IsValid() {
		*result = append(*result, cond)
	}

	return andInternal(result, conds...)
}

func andInternal(result *condAnd, conds ...Cond) Cond {
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
		condAndPool.Put(result)
		return CondEmpty
	case 1:
		// return result[0] if only one cond is valid.
		// result need to be returned to the Pool,
		// and its length resets to 0.
		retVal := (*result)[0]
		*result = (*result)[:0]
		condAndPool.Put(result)
		return retVal
	default:
		return result
	}
}

func (and *condAnd) WriteTo(w *Writer) {
	length := len(*and) - 1
	for i, cond := range *and {
		var wrap bool
		switch cond.(type) {
		case *condOr, *condExpr:
			wrap = true
			//case Eq:
			//	wrap = (len(cond.(Eq)) > 1)
			//case Neq:
			//	wrap = (len(cond.(Neq)) > 1)
		}

		if wrap {
			w.WriteByte('(')
		}

		// fmt.Println(cond)
		cond.WriteTo(w)

		if wrap {
			w.WriteByte(')')
		}

		if i != length {
			w.WriteString(" AND ")
		}
	}
}

func (and *condAnd) And(conds ...Cond) Cond {
	return andOne(and, conds...)
}

func (and *condAnd) Or(conds ...Cond) Cond {
	return orOne(and, conds...)
}

func (and *condAnd) IsValid() bool {
	return len(*and) > 1
}

func (and *condAnd) Reset() {
	if len(*and) > 0 {
		// we don't destroy cond recursively, as underlying cond can be
		// used by other sql as well, which can cause a complex dependence
		// graph. One can use the set to track destroyed conds. However, it
		// causes the overhead significantly. In addition, a destroyed cond
		// may be still referenced by other condition. Then, we have to use
		// a counter to track the number referencer which further impacts
		// the performance. Thus, the simplest way is to let user to call 
		// Destroy() or Reset to free the cond explicitly.
		//
		// No need to worry about the dirty data in the array as it is just
		// a reference and will be overwritten next time.
		*and = (*and)[:0]
	}
}

func (and *condAnd) Destroy() {
	and.Reset()
	condAndPool.Put(and)
}
