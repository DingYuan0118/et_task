package util

import (
	"bytes"
)

type File struct {
	name   string
	Buffer *bytes.Buffer
}

func NewFile(name string) *File {
	return &File{
		name:   name,
		Buffer: &bytes.Buffer{},
	}
}

func (f *File) Write(chunk []byte) error {
	_, err := f.Buffer.Write(chunk)

	return err
}