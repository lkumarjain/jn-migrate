package sql

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/lkumarjain/jn-migrate/store"
)

func (w *writer) Initialize() error {
	store.Logger.Printf("Writer Configuration : %+v", w.config)
	w.updateIdentifier()
	db, err := connection(w.config)
	if err != nil {
		return err
	}
	w.db = db
	w.createSchema(db)

	err = w.createTable(db)
	if err != nil {
		return err
	}

	return w.createStatement(db)
}

func (w *writer) Write(record store.Row) (bool, error) {
	_, err := w.statement.Exec(w.extractColumnValue(record)...)
	return err == nil, err
}

func (w *writer) Flush() error {
	err := w.statement.Close()
	if err != nil {
		return err
	}

	return w.transaction.Commit()
}

func (w *writer) updateIdentifier() {
	w.schema = toSQL(w.config.Schema)
	w.table = toSQL(w.config.Table)
	w.columnTypes = make([]string, len(w.config.Columns))
	w.columnSpecifiers = make([]string, len(w.config.Columns))
	w.columns = make([]string, len(w.config.Columns))
	for i, col := range w.config.Columns {
		column := toSQL(col)
		w.columnTypes[i] = fmt.Sprintf("%s TEXT", column)
		w.columnSpecifiers[i] = "$" + strconv.Itoa(i+1)
		w.columns[i] = column
	}
}

//createSchema : Tries to create the schema and ignores failures
func (w *writer) createSchema(db *sql.DB) {
	store.Logger.Printf("Creating Schema :: %s", w.schema)
	createSchema, err := db.Prepare(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", w.schema))
	if err == nil {
		createSchema.Exec() //nolint
	}
}

//createTable : Create a table with TEXT columns, if does not exists
func (w *writer) createTable(db *sql.DB) error {
	store.Logger.Printf("Creating Table :: %s", w.table)
	columnDefinitions := strings.Join(w.columnTypes, ",")
	fullyQualifiedTable := fmt.Sprintf("%s.%s", w.schema, w.table)
	tableSchema := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", fullyQualifiedTable, columnDefinitions)
	store.Logger.Printf("Creating Table Using :: %s", tableSchema)
	statement, err := db.Prepare(tableSchema)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	return err
}

//createStatement : Create an insert statement by openning a transaction
func (w *writer) createStatement(db *sql.DB) error {
	var err error
	w.transaction, err = db.Begin()
	if err != nil {
		return err
	}

	fullyQualifiedTable := fmt.Sprintf("%s.%s", w.schema, w.table)
	values := strings.Join(w.columnSpecifiers, ",")
	columns := strings.Join(w.columns, ",")
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", fullyQualifiedTable, columns, values) //nolint
	w.statement, err = w.transaction.Prepare(query)
	return err
}

//extractColumnValue : Extract column values from a row
func (w *writer) extractColumnValue(row store.Row) []interface{} {
	size := len(w.config.Columns)
	values := make([]interface{}, size)
	for index := 0; index < size; index++ {
		values[index] = row.GetColumn(w.config.Columns[index]).Value
	}
	return values
}
