package main

import (
	"fmt"
	"time"

	"github.com/lkumarjain/jn-migrate/store"
	"github.com/lkumarjain/jn-migrate/store/sql"
	"github.com/lkumarjain/jn-migrate/store/sql/postgres"
)

func main() {
	w := sql.Writer(sql.Config{
		ConnectionString:      "postgres://postgres:postgres@localhost:5432/postgres",
		Columns:               []string{"test", "testint"},
		MaxParallelConnection: 3,
		Dialect:               postgres.Dialect,
		Schema:                "public",
		Table:                 "test1",
	})
	err := w.Initialize()
	if err != nil {
		fmt.Println("Initialize ", err)
	}
	for index := 20; index < 30; index++ {
		_, err = w.Write(store.Row{Columns: []store.Column{
			store.Column{
				Name:  "test",
				Value: fmt.Sprintf("test-%d", index),
			},
			store.Column{
				Name:  "testint",
				Value: fmt.Sprintf("%d", index),
			},
		}})
		if err != nil {
			fmt.Println("Write ", err)
		}
	}
	err = w.Flush()
	if err != nil {
		fmt.Println("Flush ", err)
	}
	time.Sleep(time.Second)
}
