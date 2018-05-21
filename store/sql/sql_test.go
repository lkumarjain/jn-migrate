package sql

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		{
			"Config Instance",
			&Config{
				MaxParallelConnection: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriter(t *testing.T) {
	got := Writer(Config{})
	_, ok := got.(*writer)
	if !ok {
		t.Errorf("Writer() = %v, want %v", got, "writer")
	}
}

func Test_toSQL(t *testing.T) {
	type args struct {
		identifier string
	}
	tests := []struct {
		name       string
		args       args
		want       string
		exactMatch bool
	}{
		{
			name:       "Blank-String",
			args:       args{""},
			want:       "column_",
			exactMatch: false,
		},
		{
			name:       "Capital-Identifier",
			args:       args{"NAME"},
			want:       "name",
			exactMatch: true,
		},
		{
			name:       "Spacial-Charactor-Identifier",
			args:       args{"Column/1#2"},
			want:       "column1_2",
			exactMatch: true,
		},
		{
			name:       "Number-Identifier",
			args:       args{"012Column/1#2"},
			want:       "_012column1_2",
			exactMatch: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toSQL(tt.args.identifier)
			if (tt.exactMatch && got != tt.want) && (!tt.exactMatch && !strings.Contains(got, tt.want)) {
				t.Errorf("toSQL() = %v, want %v", got, tt.want)
			}
		})
	}
}
