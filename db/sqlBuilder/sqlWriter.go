package sqlBuilder

import (
	"io"
	"strings"
)

// Writer defines the interface
type Writer interface {
	io.Writer
	Append(...interface{})
}

var _ Writer = NewWriter()

// BytesWriter implments Writer and save SQL in bytes.Buffer
type BytesWriter struct {
	*strings.Builder
	args     []interface{}
	bulkArgs [][]interface{}
}

// NewWriter creates a new string writer
func NewWriter() *BytesWriter {
	w := &BytesWriter{
		Builder:  &strings.Builder{},
		args:     []interface{}{},
		bulkArgs: [][]interface{}{},
	}
	w.Grow(10)
	return w
}

// Append appends args to Writer
func (w *BytesWriter) Append(args ...interface{}) {
	w.args = append(w.args, args...)
}

func (w *BytesWriter) AppendBulk(args ...[]interface{}) {
	w.bulkArgs = append(w.bulkArgs, args...)
}

// Args returns args
func (w *BytesWriter) Args() []interface{} {
	return w.args
}

// Args returns args
func (w *BytesWriter) BulArgs() [][]interface{} {
	return w.bulkArgs
}
