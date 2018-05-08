package sql

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/kennygrant/sanitize"
	"github.com/lkumarjain/jn-migrate/store"
)

//Config is a struct to hold all postgre configuration
type Config struct {
	// ConnectionString : When establishing a connection you are expected to supply a connection string containing zero or more parameters.
	// A subset of the connection parameters supported are :
	// dbname - The name of the database to connect to
	// user - The user to sign in as
	// password - The user's password
	// host - The host to connect to. Values that start with / are for unix
	// domain sockets. (default is localhost)
	// port - The port to bind to. (default is 5432)
	// sslmode - Whether or not to use SSL (default is require, this is not the default for libpq)
	// fallback_application_name - An application_name to fall back to if one isn't provided.
	// connect_timeout - Maximum wait for connection, in seconds. Zero or not specified means wait indefinitely.
	// sslcert - Cert file location. The file must contain PEM encoded data.
	// sslkey - Key file location. The file must contain PEM encoded data.
	// sslrootcert - The location of the root certificate file. The file must contain PEM encoded data.
	ConnectionString      string
	Dialect               string
	Schema                string
	Table                 string
	Columns               []string
	MaxParallelConnection int
}

var replacements = map[string]string{
	" ": "_",
	"/": "_",
	".": "_",
	":": "_",
	";": "_",
	"|": "_",
	"-": "_",
	",": "_",
	"#": "_",

	"[":  "",
	"]":  "",
	"{":  "",
	"}":  "",
	"(":  "",
	")":  "",
	"?":  "",
	"!":  "",
	"$":  "",
	"%":  "",
	"*":  "",
	"\"": "",
}

//toSQL is a function to convert identifier name in SQL format
func toSQL(identifier string) string {
	str := sanitize.BaseName(identifier)
	str = strings.ToLower(identifier)
	str = strings.TrimSpace(str)

	for oldString, newString := range replacements {
		str = strings.Replace(str, oldString, newString, -1)
	}

	if len(str) == 0 {
		return fmt.Sprintf("column_%d", rand.Intn(10000))
	}

	firstLetter := string(str[0])
	if _, err := strconv.Atoi(firstLetter); err == nil {
		str = "_" + str
	}
	return str
}

type writer struct {
	config           Config
	openConnection   func(driverName, dataSourceName string) (*sql.DB, error)
	schema           string
	table            string
	columns          []string
	columnTypes      []string
	columnSpecifiers []string
	statement        *sql.Stmt
	transaction      *sql.Tx
}

//Writer returns a sql writer
func Writer(config Config) store.Writer {
	return &writer{config: config, openConnection: sql.Open}
}
