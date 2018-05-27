package sql

import (
	"database/sql"
	"testing"

	"github.com/lkumarjain/jn-migrate/store"
)

func Test_reader_Read(t *testing.T) {
	type fields struct {
		config    Config
		schema    string
		table     string
		columns   []string
		statement *sql.Stmt
		db        *sql.DB
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
			name:   "Initialization Error",
			fields: fields{config: Config{}},
			args: args{record: func(row store.Row) {
				if row.Error == nil {
					panic("Expected Error")
				}
			}},
			hasPanic: true,
		},
		{
			name: "Initialization Success",
			fields: fields{
				config: Config{ConnectionString: "fakedb", Dialect: "fakedb", Columns: []string{"h"}},
			},
			args: args{record: func(row store.Row) {
				if row.Error == nil {
					panic("Expected Error")
				}
			}},
			hasPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &reader{
				config:    tt.fields.config,
				schema:    tt.fields.schema,
				table:     tt.fields.table,
				columns:   tt.fields.columns,
				statement: tt.fields.statement,
				db:        tt.fields.db,
			}
			if tt.hasPanic {
				defer func() {
					r := recover()
					if r == nil {
						t.Errorf("reader.Read() expected panic got nil")
					}
				}()
			}
			r.Read(tt.args.record)
		})
	}
}
