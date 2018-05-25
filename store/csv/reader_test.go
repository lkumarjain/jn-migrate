package csv

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/lkumarjain/jn-migrate/store"
)

func Test_reader_Read(t *testing.T) {
	type fields struct {
		cfg      Config
		openFile func(name string) (*os.File, error)
		reader   func(r io.Reader) *csv.Reader
	}
	type args struct {
		record store.Record
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		hasPanic bool
	}{
		{
			name: "1-Open-Error",
			fields: fields{cfg: Config{},
				openFile: func(name string) (*os.File, error) { return nil, errors.New("Error") },
				reader:   func(r io.Reader) *csv.Reader { return csv.NewReader(strings.NewReader("in")) },
			},
			args:     args{record: func(row store.Row) {}},
			hasPanic: true,
		},
		{
			name: "2-Open-Success",
			fields: fields{cfg: Config{Header: []string{"h"}},
				openFile: func(name string) (*os.File, error) { return &os.File{}, nil },
				reader: func(r io.Reader) *csv.Reader {
					rd := csv.NewReader(strings.NewReader("in"))
					rd.FieldsPerRecord = 20
					return rd
				},
			},
			args: args{record: func(row store.Row) {
				if row.Error == nil {
					panic("Expected Error")
				}
			}},
			hasPanic: false,
		},
		{
			name: "3-Open-Success-Has-Header",
			fields: fields{cfg: Config{HasHeader: true, Header: []string{"h"}},
				openFile: func(name string) (*os.File, error) { return &os.File{}, nil },
				reader: func(r io.Reader) *csv.Reader {
					rd := csv.NewReader(strings.NewReader("in"))
					rd.FieldsPerRecord = 20
					return rd
				},
			},
			args: args{record: func(row store.Row) {
				if row.Error == nil {
					panic("Expected Error")
				}
			}},
			hasPanic: false,
		},
		{
			name: "4-Success",
			fields: fields{cfg: Config{Header: []string{"h"}},
				openFile: func(name string) (*os.File, error) { return &os.File{}, nil },
				reader:   func(r io.Reader) *csv.Reader { return csv.NewReader(strings.NewReader("in")) },
			},
			args: args{record: func(row store.Row) {
				if row.Error != nil {
					panic("Expected Error")
				}
			}},
			hasPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.hasPanic {
				defer func() {
					r := recover()
					if r == nil {
						t.Errorf("reader.Read() expected panic got nil")
					}
				}()
			}
			r := reader{
				cfg:      tt.fields.cfg,
				openFile: tt.fields.openFile,
				reader:   tt.fields.reader,
			}
			r.Read(tt.args.record)
		})
	}
}

func Test_reader_readValue(t *testing.T) {
	type fields struct {
		cfg Config
	}
	type args struct {
		result []string
		index  int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "1",
			fields: fields{},
			args:   args{result: []string{"1", "2", "3"}, index: 1},
			want:   "2",
		},
		{
			name:   "2",
			fields: fields{},
			args:   args{result: []string{"1", "2", "3"}, index: 2},
			want:   "3",
		},
		{
			name:   "3",
			fields: fields{},
			args:   args{result: []string{"1", "2", "3"}, index: -1},
			want:   "",
		},
		{
			name:   "4",
			fields: fields{},
			args:   args{result: []string{"1", "2", "3"}, index: 4},
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader{
				cfg: tt.fields.cfg,
			}
			if got := r.readValue(tt.args.result, tt.args.index); got != tt.want {
				t.Errorf("reader.readValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reader_readColumns(t *testing.T) {
	type fields struct {
		cfg Config
	}
	type args struct {
		result []string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      []store.Column
		wantPanic bool
	}{
		{
			name:      "1",
			fields:    fields{},
			args:      args{},
			want:      []store.Column{store.Column{}},
			wantPanic: true,
		},
		{
			name:   "2",
			fields: fields{cfg: Config{Header: []string{"1", "2"}}},
			args:   args{result: []string{"1", "2"}},
			want: []store.Column{
				store.Column{ColumnNumber: 0, Name: "1", Value: "1"},
				store.Column{ColumnNumber: 1, Name: "2", Value: "2"},
			},
			wantPanic: false,
		},
		{
			name:   "3",
			fields: fields{cfg: Config{Header: []string{"1", "2"}}},
			args:   args{result: []string{"1"}},
			want: []store.Column{
				store.Column{ColumnNumber: 0, Name: "1", Value: "1"},
				store.Column{ColumnNumber: 1, Name: "2", Value: ""},
			},
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if r != nil && !tt.wantPanic {
					t.Errorf("reader.readColumns() = %v, want nil", r)
				}
			}()
			r := reader{
				cfg: tt.fields.cfg,
			}
			if got := r.readColumns(tt.args.result); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("reader.readColumns() = %v, want %v", got, tt.want)
			}
		})
	}
}
