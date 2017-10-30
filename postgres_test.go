package gondolier

import (
	"testing"
)

type testUser struct {
	Id      uint64 `gondolier:"type:bigint;pk;notnull"`
	Name    string `gondolier:"type:character varying(255);notnull"`
	Age     uint   `gondolier:"type:integer";notnull`
	Picture uint64 `gondolier:"type:bigint;fk:testPicture;null"`
}

type testPicture struct {
	Id       uint64 `gondolier:"type:bigint;pk;notnull"`
	FileName string `gondolier:"type:character varying(255);notnull"`
}

type testPost struct {
	Id      uint64 `gondolier:"type:bigint;pk;notnull"`
	Post    string `gondolier:"type:character varying(255);notnull"`
	User    uint64 `gondolier:"type:bigint;fk:testUser;notnull"`
	Picture uint64 `gondolier:"type:bigint;fk:testPicture;null"`
}

/*func TestPostgresMigrator(t *testing.T) {
	Use(testdb, &Postgres{})
	Model(testUser{}, testPicture{}, testPost{})
	Migrate()
}*/

func TestPostgresDropTable(t *testing.T) {
	testCleanDb()

	if _, err := testdb.Exec(`CREATE TABLE "test_user" ("id" bigint not null)`); err != nil {
		t.Fatal(err)
	}

	Use(testdb, &Postgres{})
	Drop(testUser{})

	rows, err := testdb.Query(`SELECT EXISTS (SELECT 1
	   FROM information_schema.tables 
	   WHERE table_schema = 'public'
	   AND table_name = 'test_user')`)

	if err != nil {
		t.Fatal(err)
	}

	var exists bool
	rows.Next()
	rows.Scan(&exists)

	if exists {
		t.Fatal("Table must have been dropped")
	}
}

func TestPostgresDropTableNotExists(t *testing.T) {
	testCleanDb()
	Use(testdb, &Postgres{})
	Drop(testUser{})

	rows, err := testdb.Query(`SELECT EXISTS (SELECT 1
	   FROM information_schema.tables 
	   WHERE table_schema = 'public'
	   AND table_name = 'test_user')`)

	if err != nil {
		t.Fatal(err)
	}

	var exists bool
	rows.Next()
	rows.Scan(&exists)

	if exists {
		t.Fatal("Table must have been dropped")
	}
}

func testCleanDb() {
	testdb.Exec(`DROP TABLE "test_post"`)
	testdb.Exec(`DROP TABLE "test_user"`)
	testdb.Exec(`DROP TABLE "test_picture"`)
}
