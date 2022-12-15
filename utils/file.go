package utils

import (
	"fmt"
	"io"
	"strings"
)

type NamedFileReader interface {
	Name() string
	io.ReadSeeker
}

type namedFileReader struct {
	name   string
	reader io.ReadSeeker
}

func NewNamedFileReader(reader io.ReadSeeker, name string) NamedFileReader {
	return &namedFileReader{reader: reader, name: name}
}

func (reader *namedFileReader) Name() string {
	return reader.name
}

func (reader *namedFileReader) Read(p []byte) (n int, err error) {
	return reader.reader.Read(p)
}

func (reader *namedFileReader) Seek(offset int64, whence int) (int64, error) {
	return reader.reader.Seek(offset, whence)
}

func GenerateRandFilename(reader NamedFileReader) string {
	return fmt.Sprintf("%s.%s", RandString(8), GetFileExtension(reader.Name()))
}

func GetFileExtension(filename string) string {
	if !strings.ContainsAny(filename, ".") {
		return RandFileName("", "")
	}

	splitStr := strings.Split(filename, ".")
	return splitStr[len(splitStr)-1]
}
