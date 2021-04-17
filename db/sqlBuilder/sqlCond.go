package sqlBuilder

// Cond defines an interface
type Cond interface {
	WriteTo(Writer) error
	And(...Cond) Cond
	Or(...Cond) Cond
	IsValid() bool
}

type condEmpty struct{}

var _ Cond = condEmpty{}

// NewCond creates an empty condition
func NewCond() Cond {
	return condEmpty{}
}

func (condEmpty) WriteTo(w Writer) error {
	return nil
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

func CondToSQL(cond Cond) (string, []interface{}, error) {
	if cond == nil || !cond.IsValid() {
		return "", nil, nil
	}

	w := NewWriter()
	if err := cond.WriteTo(w); err != nil {
		return "", nil, err
	}
	return w.String(), w.args, nil
}
