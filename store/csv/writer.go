package csv

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/lkumarjain/jn-migrate/store"
)

func (w *writer) Initialize() error {
	store.Logger.Printf("Writer Configuration : %+v", w.config)
	file, err := os.Create(w.config.Path)
	if err != nil {
		return fmt.Errorf("Unable to create a file at %s, : %+v", w.config.Path, err)
	}
	w.closer = file
	w.writer = csv.NewWriter(file)
	w.writer.Comma = w.config.Comma
	w.writer.UseCRLF = true
	if len(w.config.Header) != 0 {
		return w.writer.Write(w.config.Header)
	}
	return nil
}

func (w *writer) Write(record store.Row) (bool, error) {
	values, err := w.extractColumnValue(record)
	if err != nil {
		return err == nil, err
	}
	err = w.writer.Write(values)
	if err != nil {
		return err == nil, fmt.Errorf("Unable to write column : %+v", err)
	}
	return false, nil
}

func (w *writer) Flush() error {
	w.writer.Flush()
	return w.closer.Close()
}

//extractColumnValue : Extract column values from a row
func (w *writer) extractColumnValue(row store.Row) ([]string, error) {
	size := len(w.config.Header)
	if size == 0 {
		size = len(row.Columns)
		w.config.Header = make([]string, size)
		values := make([]string, size)
		for index := 0; index < size; index++ {
			column := row.Columns[index]
			w.config.Header[index] = column.Name
			values[index] = column.Value
		}
		err := w.writer.Write(w.config.Header)
		if err != nil {
			return values, fmt.Errorf("Unable to write header : %+v", err)
		}
		return values, nil
	}
	values := make([]string, size)
	for index := 0; index < size; index++ {
		values[index] = row.GetColumn(w.config.Header[index]).Value
	}
	return values, nil
}
