package gondolier

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
	"testing"
)

var (
	testdb *sql.DB
)

func TestMain(m *testing.M) {
	// connect to database
	var err error
	testdb, err = sql.Open("postgres", testGetDbString())

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

func testGetDbString() string {
	return "host=" + os.Getenv("TEST_PG_HOST") +
		" port=" + os.Getenv("TEST_PG_PORT") +
		" user=" + os.Getenv("TEST_PG_USER") +
		" password=" + os.Getenv("TEST_PG_PASSWORD") +
		" dbname=" + os.Getenv("TEST_PG_DB") +
		" sslmode=disable"
}
