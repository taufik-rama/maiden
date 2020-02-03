package writer

import (
	"fmt"
	"os"
)

// Writer ...
type Writer struct {
	identlevel int
	file       *os.File
}

// New ...
func New(filename string) (Writer, error) {
	file, err := os.Create(filename)
	if err != nil {
		return Writer{}, err
	}
	return Writer{
		file: file,
	}, nil
}

// DecrIndentLevel ...
func (w *Writer) DecrIndentLevel() {
	w.identlevel--
}

// IncrIndentLevel ...
func (w *Writer) IncrIndentLevel() {
	w.identlevel++
}

// WriteRaw ...
func (w Writer) WriteRaw(format string, v ...interface{}) error {
	_, err := w.file.WriteString(fmt.Sprintf(format, v...))
	return err
}

// Write ...
func (w Writer) Write(format string, v ...interface{}) error {
	ident := ""
	for i := 0; i < w.identlevel; i++ {
		ident += "\t"
	}
	return w.WriteRaw((ident + format + "\n"), v...)
}

// WriteEmptyLine ...
func (w Writer) WriteEmptyLine() error {
	_, err := w.file.WriteString("\n")
	return err
}

// Close ...
func (w Writer) Close() error {
	return w.file.Close()
}
