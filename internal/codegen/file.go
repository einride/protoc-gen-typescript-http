package codegen

import (
	"bytes"
	"fmt"
)

type File struct {
	buf bytes.Buffer
}

func (f *File) P(v ...interface{}) {
	for _, x := range v {
		fmt.Fprint(&f.buf, x)
	}
	fmt.Fprintln(&f.buf)
}

func (f *File) Content() []byte {
	return f.buf.Bytes()
}
