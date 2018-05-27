package sql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/lkumarjain/jn-migrate/store"
)

func (r *reader) Read(record store.Record) {
	err := r.initialize()
	if err != nil {
		panic(err)
	}

	rows, err := r.statement.Query()
	if err != nil {
		panic(fmt.Errorf("Failed to execute Query %+v", err))
	}
	defer rows.Close()
	defer r.statement.Close()
	defer r.db.Close()
	columnNames, err := rows.Columns()
	if err != nil {
		panic(fmt.Errorf("Failed to find column %+v", err))
	}

	count := len(columnNames)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	rowNum := int64(0)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		err = rows.Scan(valuePtrs...)
		if err != nil {
			record(store.Row{RowNumber: rowNum, Error: err})
			continue
		}
		columns := make([]store.Column, count)
		for i := 0; i < count; i++ {
			columns[i] = store.Column{ColumnNumber: i, Name: columnNames[i], Value: fmt.Sprintf("%v", values[i])}
		}
		record(store.Row{RowNumber: rowNum, Columns: columns})
	}
}

func (r *reader) initialize() error {
	store.Logger.Printf("Reader Configuration : %+v", r.config)
	r.updateIdentifier()
	db, err := connection(r.config)
	if err != nil {
		return err
	}
	r.db = db
	return r.createStatement(db)
}

//createStatement : Create an insert statement by openning a transaction
func (r *reader) createStatement(db *sql.DB) error {
	var err error
	fullyQualifiedTable := fmt.Sprintf("%s.%s", r.schema, r.table)
	columns := strings.Join(r.columns, ",")
	query := fmt.Sprintf("SELECT %s FROM %s", columns, fullyQualifiedTable) //nolint
	r.statement, err = db.Prepare(query)
	return err
}

func (r *reader) updateIdentifier() {
	r.schema = toSQL(r.config.Schema)
	r.table = toSQL(r.config.Table)
	r.columns = make([]string, len(r.config.Columns))
	for i, col := range r.config.Columns {
		column := toSQL(col)
		r.columns[i] = column
	}
}
