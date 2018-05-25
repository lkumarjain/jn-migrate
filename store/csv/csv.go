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
	//CSV headers; If blank will be populated automatically, by reading CSV or Row
	Header []string
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
	openFile func(name string) (*os.File, error)
	reader   func(r io.Reader) *csv.Reader
}

//Reader returns a CSV file reader
func Reader(cfg Config) store.Reader {
	return reader{cfg: cfg, openFile: os.Open, reader: csv.NewReader}
}

type writer struct {
	config Config
	writer *csv.Writer
	closer io.Closer
}

//Writer returns a CSV writer
func Writer(config Config) store.Writer {
	return &writer{config: config}
}

//InitiailzedWriter returns a Initiailzed CSV writer
func InitiailzedWriter(config Config) (store.Writer, error) {
	w := &writer{config: config}
	err := w.Initialize()
	return w, err
}
