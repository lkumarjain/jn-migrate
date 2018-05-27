package main

import (
	"fmt"
	"time"

	"github.com/lkumarjain/jn-migrate/store"
	"github.com/lkumarjain/jn-migrate/store/sql"
	"github.com/lkumarjain/jn-migrate/store/sql/postgres"
)

func main() {
	r := sql.Reader(sql.Config{
		ConnectionString:      "postgres://postgres:postgres@localhost:5432/postgres",
		Columns:               []string{"test", "testint"},
		MaxParallelConnection: 3,
		Dialect:               postgres.Dialect,
		Schema:                "import",
		Table:                 "test",
	})

	r.Read(process)
	time.Sleep(time.Second)
}

func process(row store.Row) {
	fmt.Printf("Record : %v, Error :: %v\n", row.Columns, row.Error)
	if row.Error != nil {
		//fmt.Printf("Record : %+v, Error :: %v\n", row.RowNumber, row.Error)
		fmt.Printf("Record : %+v\n", row)
	}
}
