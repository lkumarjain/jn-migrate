package sql

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/codegangsta/cli"
	"github.com/lkumarjain/jn-migrate/store"
	"github.com/lkumarjain/jn-migrate/store/csv"
)

const (
	//CSV Flags
	skipHeader = "skip-header"
	fields     = "fields"
	delimiter  = "delimiter"
	trimSpace  = "trim-space"
	source     = "source"
	comment    = "comment"

	//SQL flags
	connectionURL = "connection-url"
	dbName        = "dbname"
	schema        = "schema"
	table         = "table"
	poolSize      = "poolsize"
)

//Command for CSV to SQL persistance
var Command = cli.Command{
	Name:  "csvtosql",
	Usage: "Import csv into database",
	Flags: []cli.Flag{
		//CSV Flags
		cli.BoolFlag{
			Name:  fmt.Sprintf("%s, %s", skipHeader, "sh"),
			Usage: "skip header row",
		},
		cli.StringFlag{
			Name:  fmt.Sprintf("%s, %s", fields, "f"),
			Usage: "comma separated field names if no header row",
		},
		cli.StringFlag{
			Name:  fmt.Sprintf("%s, %s", delimiter, "d"),
			Value: ",",
			Usage: "field delimiter",
		},
		cli.BoolFlag{
			Name:  fmt.Sprintf("%s, %s", trimSpace, "ts"),
			Usage: "leading white space in a field is ignored",
		},
		cli.StringFlag{
			Name:  fmt.Sprintf("%s, %s", source, "s"),
			Value: "input.csv",
			Usage: "CSV file path",
		},
		cli.StringFlag{
			Name:  fmt.Sprintf("%s, %s", comment, "co"),
			Value: "#",
			Usage: "Comment character without preceding whitespace are ignored",
		},
		//SQL Flags
		cli.StringFlag{
			Name:   fmt.Sprintf("%s, %s", connectionURL, "c"),
			Usage:  "When establishing a connection you are expected to supply a connection string containing zero or more parameters",
			EnvVar: "DB_CONNECTION_URL",
		},
		cli.StringFlag{
			Name:   fmt.Sprintf("%s, %s", dbName, "db"),
			Value:  "postgres",
			Usage:  "database to connect to",
			EnvVar: "DB_NAME",
		},
		cli.StringFlag{
			Name:   fmt.Sprintf("%s, %s", schema, "sc"),
			Value:  "import",
			Usage:  "database schema",
			EnvVar: "DB_SCHEMA",
		},
		cli.StringFlag{
			Name:   fmt.Sprintf("%s, %s", table, "t"),
			Usage:  "destination table",
			EnvVar: "DB_TABLE",
		},
		cli.IntFlag{
			Name:  fmt.Sprintf("%s, %s", poolSize, "p"),
			Value: 10,
			Usage: "Maximum Parallel Connection",
		},
	},
	Action: action,
}

func action(c *cli.Context) {
	cli.CommandHelpTemplate = strings.Replace(cli.CommandHelpTemplate, "[arguments...]", "<csv-file>", -1)
	conf := csv.NewConfig()
	conf.Path = c.String(source)
	conf.Comma, _ = utf8.DecodeRuneInString(c.String(delimiter))
	conf.Comment, _ = utf8.DecodeRuneInString(c.String(comment))
	conf.TrimLeadingSpace = c.Bool(trimSpace)
	fields := c.String(fields)
	conf.HasHeader = fields == ""
	r := csv.Reader(*conf)
	r.Read(process)
}

func process(row store.Row) {
	if row.Error != nil {
		//fmt.Printf("Record : %+v, Error :: %v\n", row.RowNumber, row.Error)
		fmt.Printf("Record : %+v\n", row)
	}
}
