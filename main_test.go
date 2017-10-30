package gondolier

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
	"testing"
)

const (
	db_conn = "host=localhost port=5432 user=postgres password=postgres dbname=gondolier sslmode=disable"
)

var (
	testdb *sql.DB
)

func TestMain(m *testing.M) {
	// connect to database
	var err error
	testdb, err = sql.Open("postgres", db_conn)

	if err != nil {
		panic(err)
	}

	if err := testdb.Ping(); err != nil {
		panic(err)
	}

	// run
	code := m.Run()
	os.Exit(code)
}
