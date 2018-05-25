package sql

import (
	"database/sql"
	"testing"

	_ "github.com/andizzle/go-fakedb"
	"github.com/lkumarjain/jn-migrate/store"
)

func Test_writer_Initialize(t *testing.T) {
	type fields struct {
		config           Config
		schema           string
		table            string
		columns          []string
		columnTypes      []string
		columnSpecifiers []string
		statement        *sql.Stmt
		transaction      *sql.Tx
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "Blank-Connection-String-Failure",
			fields:  fields{config: Config{Columns: []string{"1", "2"}}},
			wantErr: true,
		},
		{
			name: "Open-Connection-Failure",
			fields: fields{
				config: Config{ConnectionString: "1", Columns: []string{"1", "2"}},
			},
			wantErr: true,
		},
		{
			name: "Open-Connection-Success",
			fields: fields{
				config: Config{
					ConnectionString: "1",
					Columns:          []string{"1", "2"},
					Dialect:          "fakedb",
					Table:            "fakeTable",
					Schema:           "fakeSchema",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &writer{
				config:           tt.fields.config,
				schema:           tt.fields.schema,
				table:            tt.fields.table,
				columns:          tt.fields.columns,
				columnTypes:      tt.fields.columnTypes,
				columnSpecifiers: tt.fields.columnSpecifiers,
				statement:        tt.fields.statement,
				transaction:      tt.fields.transaction,
			}
			if err := w.Initialize(); (err != nil) != tt.wantErr {
				t.Errorf("writer.Initialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_writer_Write(t *testing.T) {
	type fields struct {
		config           Config
		schema           string
		table            string
		columns          []string
		columnTypes      []string
		columnSpecifiers []string
		statement        *sql.Stmt
		transaction      *sql.Tx
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
			name: "statement-execute",
			fields: fields{config: Config{
				ConnectionString: "1",
				Columns:          []string{"1", "2"},
				Dialect:          "fakedb",
				Table:            "fakeTable",
				Schema:           "fakeSchema",
			},
				statement: &sql.Stmt{},
			},
			args: args{record: store.Row{
				Columns: []store.Column{
					store.Column{Name: "1", Value: "1"},
				}},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &writer{
				config:           tt.fields.config,
				schema:           tt.fields.schema,
				table:            tt.fields.table,
				columns:          tt.fields.columns,
				columnTypes:      tt.fields.columnTypes,
				columnSpecifiers: tt.fields.columnSpecifiers,
				statement:        tt.fields.statement,
				transaction:      tt.fields.transaction,
			}
			w.Initialize()
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

func Test_writer_Flush(t *testing.T) {
	type fields struct {
		config           Config
		schema           string
		table            string
		columns          []string
		columnTypes      []string
		columnSpecifiers []string
		statement        *sql.Stmt
		transaction      *sql.Tx
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "statement-execute",
			fields: fields{config: Config{
				ConnectionString: "1",
				Columns:          []string{"1", "2"},
				Dialect:          "fakedb",
				Table:            "fakeTable",
				Schema:           "fakeSchema",
			},
				statement: &sql.Stmt{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &writer{
				config:           tt.fields.config,
				schema:           tt.fields.schema,
				table:            tt.fields.table,
				columns:          tt.fields.columns,
				columnTypes:      tt.fields.columnTypes,
				columnSpecifiers: tt.fields.columnSpecifiers,
				statement:        tt.fields.statement,
				transaction:      tt.fields.transaction,
			}
			w.Initialize()
			if err := w.Flush(); (err != nil) != tt.wantErr {
				t.Errorf("writer.Flush() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
