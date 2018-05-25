package main

import (
	"fmt"

	"github.com/lkumarjain/jn-migrate/store"
	"github.com/lkumarjain/jn-migrate/store/csv"
)

func write() {
	cfg := csv.NewConfig()
	cfg.Path = "./result.csv"
	cfg.Comma = ','
	cfg.Header = []string{"Name-1", "Name-4"}
	w, err := csv.InitiailzedWriter(*cfg)
	if err != nil {
		fmt.Printf("InitiailzedWriter cfg : %+v, %v", cfg, err)
	}
	for index := int64(0); index < 10; index++ {
		row := store.Row{
			RowNumber: index,
			Columns: []store.Column{
				store.Column{ColumnNumber: 1, Name: "Name-1", Value: fmt.Sprintf("Value-1-%d", index)},
				store.Column{ColumnNumber: 2, Name: "Name-2", Value: fmt.Sprintf("Value-2-%d", index)},
				store.Column{ColumnNumber: 3, Name: "Name-3", Value: fmt.Sprintf("Value-3-%d", index)},
				store.Column{ColumnNumber: 4, Name: "Name-4", Value: fmt.Sprintf("Value-4-%d", index)},
			},
		}
		_, err = w.Write(row)
		if err != nil {
			fmt.Printf("Write Row : %+v, %v", row, err)
		}
	}

	err = w.Flush()
	if err != nil {
		fmt.Printf("Flush %+v", err)
	}
}
