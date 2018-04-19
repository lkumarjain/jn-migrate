package store

//Row is a struct to hold all csv values
type Row struct {
	RowNumber int64
	Columns   []Column
	Error     error
}

//Column is a struct to hold all csv values
type Column struct {
	ColumnNumber int
	Name         string
	Value        string
}

//Record is a function called by Reader
type Record func(row Row)

//Reader is an interface to read all the records from a database
type Reader interface {
	//Read is a function to read all the records from database
	//Record is a callback function of record to be returned to caller
	Read(record Record)
}
