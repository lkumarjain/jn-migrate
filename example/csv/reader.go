package main

import (
	"fmt"

	"github.com/lkumarjain/jn-migrate/store"
	"github.com/lkumarjain/jn-migrate/store/csv"
)

func main() {
	cfg := csv.NewConfig()
	cfg.Path = "./data.cvs"
	cfg.Comma = ';'
	//cfg.Path:             "/home/lokesh/Downloads/DC-43251/agent_provisioning.csv",
	r := csv.Reader(*cfg)
	r.Read(process)
}

func process(row store.Row) {
	if row.Error != nil {
		//fmt.Printf("Record : %+v, Error :: %v\n", row.RowNumber, row.Error)
		fmt.Printf("Record : %+v\n", row)
	}
}
