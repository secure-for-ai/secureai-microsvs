package sqlBuilderV3

import _ "unsafe"

// Cond defines an interface
type Cond interface {
	WriteTo(*Writer)
	And(...Cond) Cond
	Or(...Cond) Cond
	IsValid() bool
	Reset()
	Destroy()
}

type condEmpty struct{}

var _ Cond = condEmpty{}
var CondEmpty = condEmpty{}

// NewCond creates an empty condition
func NewCond() Cond {
	return condEmpty{}
}

func (condEmpty) WriteTo(w *Writer) {
}

func (condEmpty) And(conds ...Cond) Cond {
	return And(conds...)
}

func (condEmpty) Or(conds ...Cond) Cond {
	return Or(conds...)
}

func (condEmpty) IsValid() bool {
	return false
}

func (condEmpty) Reset() {
}

func (condEmpty) Destroy() {
}

func CondToSQL(cond Cond, w *Writer) (string, []any, error) {
	if cond == nil || !cond.IsValid() {
		return "", []any{}, nil
	}

	w.Reset()
	cond.WriteTo(w)
	return w.String(), w.Args(), nil
}
