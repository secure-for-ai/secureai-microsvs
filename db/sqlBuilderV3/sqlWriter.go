package sqlBuilderV3

import (
	"bytes"

	"github.com/secure-for-ai/secureai-microsvs/util"
)

// Writer defines the interface
// type Writer interface {
// 	io.Writer
// 	io.StringWriter
// 	io.ByteWriter
// 	Append(...interface{})
// 	String() string
// 	Args() []interface{}
// 	Reset()
// }

func writerJoin(w *Writer, elems []string, sep byte) {
	w.WriteString(elems[0])
	for _, s := range elems[1:] {
		w.WriteByte(sep)
		w.WriteString(s)
	}
}

// var _ Writer = NewWriter()

// Writer implments Writer and save SQL in bytes.Buffer
type Writer struct {
	*bytes.Buffer
	args     []interface{}
	bulkArgs [][]interface{}
}

// NewWriter creates a new string writer
func NewWriter() *Writer {
	w := &Writer{
		Buffer:   &bytes.Buffer{},
		args:     []interface{}{},
		bulkArgs: [][]interface{}{},
	}
	w.Grow(128)
	return w
}

// Append appends args to Writer
func (w *Writer) Append(args ...interface{}) {
	w.args = append(w.args, args...)
}

func (w *Writer) AppendBulk(args ...[]interface{}) {
	w.bulkArgs = append(w.bulkArgs, args...)
}

// Args returns args
func (w *Writer) Args() []interface{} {
	return w.args
}

// Args returns args
func (w *Writer) BulkArgs() [][]interface{} {
	return w.bulkArgs
}

func (w *Writer) String() string {
	return util.FastBytesToString(w.Bytes())
}

func (w *Writer) Reset() {
	w.Buffer.Reset()
	w.args = w.args[:0]
	w.bulkArgs = w.bulkArgs[:0]
}
