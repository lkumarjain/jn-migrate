package postgres

import (
	_ "github.com/lib/pq" //Loading postgre driver
)

//Dialect is a database name used for registration
const Dialect = "postgres"
