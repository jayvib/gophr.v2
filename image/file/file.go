package file

import "io"

type File interface {
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Closer
}

type Metadata struct {
	Filename string
	Description string
}