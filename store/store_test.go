package store

import (
	"reflect"
	"testing"
)

func TestRow_GetColumn(t *testing.T) {
	type fields struct {
		RowNumber   int64
		Columns     []Column
		Error       error
		name2Column map[string]Column
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Column
	}{
		{
			"Blank-Column-Name",
			fields{},
			args{},
			Column{},
		},
		{
			"Return-Value",
			fields{Columns: []Column{Column{Name: "1", Value: "2"}}},
			args{"1"},
			Column{Name: "1", Value: "2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Row{
				RowNumber:   tt.fields.RowNumber,
				Columns:     tt.fields.Columns,
				Error:       tt.fields.Error,
				name2Column: tt.fields.name2Column,
			}
			if got := r.GetColumn(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Row.GetColumn() = %v, want %v", got, tt.want)
			}
		})
	}
}
