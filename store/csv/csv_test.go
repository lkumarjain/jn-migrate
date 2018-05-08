package csv

import (
	"testing"
)

func TestGetReader(t *testing.T) {
	got := Reader(Config{})
	_, ok := got.(reader)
	if !ok {
		t.Errorf("Reader() = %v, want %v", got, reader{})
	}
}

func TestGetReaderWithHeader(t *testing.T) {
	got := ReaderWithHeader(Config{}, []string{"1", "2"})
	_, ok := got.(reader)
	if !ok {
		t.Errorf("ReaderWithHeader() = %v, want %v", got, reader{})
	}
}
