package main

import (
	"fmt"

	"github.com/lkumarjain/jn-migrate/store"
	"github.com/lkumarjain/jn-migrate/store/csv"
)

func main() {
	r := csv.Reader(csv.Config{
		Path:  "./data.cvs",
		Comma: ';',
		// Path:             "/home/lokesh/Downloads/DC-43251/agent_provisioning.csv",
		// Comma:            ',',
		Comment:          '#',
		HasHeader:        true,
		TrimLeadingSpace: true,
	})

	r.Read(process)
}

func process(row store.Row) {
	if row.Error != nil {
		//fmt.Printf("Record : %+v, Error :: %v\n", row.RowNumber, row.Error)
		fmt.Printf("Record : %+v\n", row)
	}
}
