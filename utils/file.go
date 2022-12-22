package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
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

func DownloadFile(filepath string, url string) (http.Header, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return resp.Header, err
}
