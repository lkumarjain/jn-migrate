package sql

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/lkumarjain/jn-migrate/store/sql/postgres"

	"github.com/codegangsta/cli"
	"github.com/lkumarjain/jn-migrate/store"
	"github.com/lkumarjain/jn-migrate/store/csv"
	"github.com/lkumarjain/jn-migrate/store/sql"
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

//Command for CSV to SQL persistence
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
			Value:  postgres.Dialect,
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

func action(ctx *cli.Context) {
	cli.CommandHelpTemplate = strings.Replace(cli.CommandHelpTemplate, "[arguments...]", "<csv-file>", -1)

	sCfg := sql.NewConfig()
	sCfg.ConnectionString = ctx.String(connectionURL)
	sCfg.Dialect = ctx.String(dbName)
	sCfg.Schema = ctx.String(schema)
	sCfg.Table = ctx.String(table)
	sCfg.MaxParallelConnection = ctx.Int(poolSize)
	sCfg.Columns = strings.Split(ctx.String(fields), ",")

	handler := &writeHandler{
		ignoreErrors: ctx.GlobalBool("ignore-errors"),
		writer:       sql.Writer(*sCfg),
	}

	if len(sCfg.Columns) != 0 {
		handler.initialize()
	}

	conf := csv.NewConfig()
	conf.Path = ctx.String(source)
	conf.Comma, _ = utf8.DecodeRuneInString(ctx.String(delimiter))
	conf.Comment, _ = utf8.DecodeRuneInString(ctx.String(comment))
	conf.TrimLeadingSpace = ctx.Bool(trimSpace)
	conf.HasHeader = true
	conf.Header = sCfg.Columns

	r := csv.Reader(*conf)
	r.Read(handler.process)

	err := handler.writer.Flush()
	if err != nil {
		store.Logger.Printf("Failed to write Records : %+v", err)
	}
}

type writeHandler struct {
	ignoreErrors bool
	writer       store.Writer
}

func (h *writeHandler) initialize() {
	err := h.writer.Initialize()
	if err != nil {
		panic(fmt.Errorf("Failed to initialize writer : %+v", err))
	}
}

func (h *writeHandler) process(row store.Row) {
	if row.Error != nil {
		if h.ignoreErrors {
			store.Logger.Printf("Failed to read Record : %+v", row)
		} else {
			panic(fmt.Errorf("Failed to read Record : %+v", row))
		}
		return
	}
	_, err := h.writer.Write(row)
	if err != nil {
		if h.ignoreErrors {
			store.Logger.Printf("Failed to write Record : %+v with error : %+v", row, err)
		} else {
			panic(fmt.Errorf("Failed to write Record : %+v with error : %+v", row, err))
		}
		return
	}
}
