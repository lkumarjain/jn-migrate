package csv

import (
	"fmt"
	"io"

	"github.com/lkumarjain/jn-migrate/store"
)

func (r reader) Read(record store.Record) {
	file, err := r.openFile(r.cfg.Path)
	if err != nil {
		panic(fmt.Errorf("Unable to Open File :: %s with Error :: %v", r.cfg.Path, err))
	}
	defer file.Close() //nolint
	store.Logger.Printf("Reading csv using configuration : %+v", r.cfg)
	r.read(file, record)
}

func (r reader) read(file io.Reader, record store.Record) {
	reader := r.reader(file)
	reader.Comma = r.cfg.Comma
	reader.Comment = r.cfg.Comment
	reader.TrimLeadingSpace = r.cfg.TrimLeadingSpace
	count := int64(-1)
	for {
		count = count + 1
		result, err := reader.Read()
		if err == io.EOF {
			break
		}

		if count == 0 && r.cfg.HasHeader {
			if len(r.cfg.Header) == 0 {
				r.cfg.Header = result
			}
			continue
		}
		record(store.Row{RowNumber: count, Columns: r.readColumns(result), Error: err})
	}
}

func (r reader) readColumns(result []string) []store.Column {
	size := len(r.cfg.Header)
	if size == 0 {
		panic(fmt.Errorf("Header value not available for file :: %s", r.cfg.Path))
	}

	columns := make([]store.Column, size)
	index := 0
	for ; index < size; index++ {
		columns[index] = store.Column{ColumnNumber: index, Name: r.cfg.Header[index], Value: r.readValue(result, index)}
	}
	return columns
}

func (r reader) readValue(result []string, index int) string {
	if index >= 0 && len(result) > index {
		return result[index]
	}
	return ""
}
