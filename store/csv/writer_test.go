package csv

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/lkumarjain/jn-migrate/store"
)

type writeCloser struct {
	err bool
}

func (c writeCloser) Close() error {
	if c.err {
		return errors.New("Error")
	}
	return nil
}

func (c writeCloser) Write(p []byte) (n int, err error) {
	if c.err {
		return 0, errors.New("Error")
	}
	return 0, nil
}

func Test_writer_Initialize(t *testing.T) {
	type fields struct {
		config Config
		writer *csv.Writer
		file   *os.File
	}
	tests := []struct {
		name       string
		fields     fields
		wantErr    bool
		removeFile string
	}{
		{
			name: "Create File Error",
			fields: fields{
				config: Config{Path: "//Config"},
			},
			wantErr: true,
		},
		{
			name: "Create File Success",
			fields: fields{
				config: Config{Path: "test.csv"},
			},
			wantErr:    false,
			removeFile: "test.csv",
		},
		{
			name: "Create File Success With Header",
			fields: fields{
				config: Config{Path: "test.csv", Header: []string{"1", "2"}},
			},
			wantErr:    false,
			removeFile: "test.csv",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &writer{
				config: tt.fields.config,
				writer: tt.fields.writer,
				closer: tt.fields.file,
			}
			if err := w.Initialize(); (err != nil) != tt.wantErr {
				t.Errorf("writer.Initialize() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.removeFile != "" {
				err := os.Remove(tt.removeFile)
				if err != nil {
					t.Errorf("writer.Initialize() error = %v, wantErr %v", err, false)
				}
			}
		})
	}
}

func Test_writer_Flush(t *testing.T) {
	type fields struct {
		config Config
		writer *csv.Writer
		closer io.Closer
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Flush Error",
			fields: fields{
				config: Config{},
				writer: csv.NewWriter(writeCloser{}),
				closer: writeCloser{true},
			},
			wantErr: true,
		},
		{
			name: "Flush Success",
			fields: fields{
				config: Config{},
				writer: csv.NewWriter(writeCloser{false}),
				closer: writeCloser{false},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &writer{
				config: tt.fields.config,
				writer: tt.fields.writer,
				closer: tt.fields.closer,
			}
			if err := w.Flush(); (err != nil) != tt.wantErr {
				t.Errorf("writer.Flush() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_writer_Write(t *testing.T) {
	type fields struct {
		config Config
		writer *csv.Writer
		closer io.Closer
	}
	type args struct {
		record store.Row
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Extract Column Success",
			fields: fields{
				config: Config{},
				writer: csv.NewWriter(writeCloser{true}),
				closer: writeCloser{true},
			},
			args: args{
				record: store.Row{
					Columns: []store.Column{
						store.Column{ColumnNumber: 0, Name: "1", Value: "***"},
						store.Column{ColumnNumber: 1, Name: "2", Value: "2"},
					},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Has Header Success",
			fields: fields{
				config: Config{Header: []string{"1"}},
				writer: csv.NewWriter(writeCloser{true}),
				closer: writeCloser{true},
			},
			args: args{
				record: store.Row{
					Columns: []store.Column{
						store.Column{ColumnNumber: 0, Name: "1", Value: "***"},
						store.Column{ColumnNumber: 1, Name: "2", Value: "2"},
					},
				},
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &writer{
				config: tt.fields.config,
				writer: tt.fields.writer,
				closer: tt.fields.closer,
			}
			got, err := w.Write(tt.args.record)
			if (err != nil) != tt.wantErr {
				t.Errorf("writer.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("writer.Write() = %v, want %v", got, tt.want)
			}
		})
	}
}
