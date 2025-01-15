package sqlBuilderV3

import (
	"sync"
)

var argsPool = sync.Pool{
	New: func() interface{} {
		args := make([]interface{}, 0, 4)
		return &args
	},
}

func getArgs() *[]interface{} {
	args := argsPool.Get().(*[]interface{})
	*args = (*args)[:0]
	return args
}

// Writer implments Writer and save SQL in bytes.Buffer
type Writer struct {
	*stringWriter
	args     []interface{}
	bulkArgs []*[]interface{}
}

var writerPool = sync.Pool{
	New: func() interface{} {
		w := &Writer{
			stringWriter: &stringWriter{},
			args:     make([]interface{}, 0, 4),
			bulkArgs: make([]*[]interface{}, 0, 4),
		}
		w.Grow(128)
		return w
	},
}

// NewWriter creates a new string writer
func NewWriter() *Writer {
	return writerPool.Get().(*Writer)
}

// Append appends args to Writer
func (w *Writer) Append(args ...interface{}) {
	w.args = append(w.args, args...)
}

func (w *Writer) AppendBulk(args *[]interface{}) {
	w.bulkArgs = append(w.bulkArgs, args)
}

// Args returns args
func (w *Writer) Args() []interface{} {
	return w.args
}

// Args returns args
func (w *Writer) BulkArgs() []*[]interface{} {
	return w.bulkArgs
}

func (w *Writer) Reset() {
	w.stringWriter.Reset()
	w.args = w.args[:0]
	for _, args := range w.bulkArgs {
		argsPool.Put(args)
	}
	w.bulkArgs = w.bulkArgs[:0]
}

func (w *Writer) Destroy() {
	w.Reset()
	writerPool.Put(w)
}

func (w *Writer) Join(elems []string, sep byte) {
	w.WriteString(elems[0])
	for _, s := range elems[1:] {
		w.WriteByte(sep)
		w.WriteString(s)
	}
}
