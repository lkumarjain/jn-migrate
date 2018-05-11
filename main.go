package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/lkumarjain/jn-migrate/command/sql"
)

func main() {
	app := cli.NewApp()
	app.Name = "JN-Migrate"
	app.Version = "1.0"
	app.Usage = "Import/Export from one data store to another (e.g. CSV2Database, Database2CSV, Cassandra2Database, Database2Cassandra etc)"
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

	app.Run(os.Args)
}
