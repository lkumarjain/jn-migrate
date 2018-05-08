package store

//Row is a struct to hold all csv values
type Row struct {
	RowNumber   int64
	Columns     []Column
	Error       error
	name2Column map[string]Column
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

//Writer is an interface to write all the records in the a database
type Writer interface {
	Initialize() error
	//Write is a function to write all the records in the database
	Write(record Row) (bool, error)
	Flush() error
}

//GetColumn is a function to return value for a column name
func (r Row) GetColumn(name string) Column {
	if r.name2Column == nil || len(r.name2Column) == 0 {
		r.name2Column = make(map[string]Column)
		for _, col := range r.Columns {
			r.name2Column[col.Name] = col
		}
	}
	return r.name2Column[name]
}
