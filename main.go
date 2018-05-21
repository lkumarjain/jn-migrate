package main

import (
	"log"
	"os"
	"time"

	"github.com/codegangsta/cli"
	"github.com/lkumarjain/jn-migrate/command/sql"
	"github.com/lkumarjain/jn-migrate/store"
)

func main() {
	store.Logger = log.New(os.Stdout, "[jn-migrate] ", (log.Ldate | log.Ltime | log.LUTC | log.Lshortfile))
	app := cli.NewApp()
	app.Name = "JN-Migrate"
	app.Version = "1.0"
	app.Compiled = time.Now()
	app.Copyright = "(c) 2018 lkumarjain"
	app.ArgsUsage = "[args and such]"
	app.HideHelp = false
	app.HideVersion = false
	app.Usage = "Import/Export from one data store to another"
	app.UsageText = "Import/Export from one data store to another (e.g. csvtosql, sqltocsv, cassandratosql, sqltocassandra etc)"
	app.Authors = []cli.Author{
		{
			Name:  "Lokesh Jain",
			Email: "http://lkumarjain.blogspot.net",
		},
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "ignore-errors, i",
			Usage: "halt transaction on inconsistencies",
		},
	}

	app.Commands = []cli.Command{
		sql.Command,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
