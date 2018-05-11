package csv

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/lkumarjain/jn-migrate/store"
)

//Config is a struct to hold all csv configuration
type Config struct {
	// Path CSV file path
	Path string
	// Comma is the field delimiter.
	Comma rune
	// Comment character without preceding whitespace are ignored.
	// With leading whitespace the Comment character becomes part of the
	// field, even if TrimLeadingSpace is true.
	Comment rune
	// If TrimLeadingSpace is true, leading white space in a field is ignored.
	// This is done even if the field delimiter, Comma, is white space.
	TrimLeadingSpace bool

	//HasHeader is CSV has a header? default value is true
	HasHeader bool
}

//NewConfig is a function to create configuration object
func NewConfig() *Config {
	return &Config{
		Comma:            ',',
		Comment:          '#',
		TrimLeadingSpace: true,
		HasHeader:        true,
	}
}

type reader struct {
	cfg      Config
	header   []string
	openFile func(name string) (*os.File, error)
	reader   func(r io.Reader) *csv.Reader
}

//Reader returns a CSV file reader
func Reader(cfg Config) store.Reader {
	return reader{cfg: cfg, openFile: os.Open, reader: csv.NewReader}
}

//ReaderWithHeader returns a CSV file reader with header information
func ReaderWithHeader(cfg Config, header []string) store.Reader {
	return reader{cfg: cfg, header: header, openFile: os.Open, reader: csv.NewReader}
}
