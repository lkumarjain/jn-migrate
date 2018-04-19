package csv

import (
	"testing"
)

func TestGetReader(t *testing.T) {
	got := GetReader(Config{})
	_, ok := got.(reader)
	if !ok {
		t.Errorf("GetReader() = %v, want %v", got, reader{})
	}
}

func TestGetReaderWithHeader(t *testing.T) {
	got := GetReaderWithHeader(Config{}, []string{"1", "2"})
	_, ok := got.(reader)
	if !ok {
		t.Errorf("GetReaderWithHeader() = %v, want %v", got, reader{})
	}
}
