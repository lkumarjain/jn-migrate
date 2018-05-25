package csv

import (
	"os"
	"reflect"
	"testing"
)

func TestGetReader(t *testing.T) {
	got := Reader(Config{})
	_, ok := got.(reader)
	if !ok {
		t.Errorf("Reader() = %v, want %v", got, reader{})
	}
}

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		{
			name: "Config",
			want: &Config{
				Comma:            ',',
				Comment:          '#',
				TrimLeadingSpace: true,
				HasHeader:        true,
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
		t.Errorf("Writer() = %v", got)
	}
}

func TestInitiailzedWriter(t *testing.T) {
	type args struct {
		config Config
	}
	tests := []struct {
		name       string
		args       args
		removeFile string
		wantErr    bool
	}{
		{
			name:    "Initialization Error",
			args:    args{config: Config{Path: "//Config"}},
			wantErr: true,
		},
		{
			name:       "Initialization Success",
			args:       args{config: Config{Path: "test.csv"}},
			wantErr:    false,
			removeFile: "test.csv",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InitiailzedWriter(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitiailzedWriter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			_, ok := got.(*writer)
			if !ok {
				t.Errorf("InitiailzedWriter() = %v", got)
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
