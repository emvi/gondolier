package gondolier

import (
	"testing"
)

type testUser struct {
	Id      uint64 `gondolier:"type:bigint;pk;seq:1,1,-,-,1;default:nextval(seq);notnull"`
	Name    string `gondolier:"type:character varying(255);notnull;unique"`
	Age     uint   `gondolier:"type:integer;notnull"`
	Picture uint64 `gondolier:"type:bigint;fk:testPicture;null"`
}

type testPicture struct {
	Id       uint64 `gondolier:"type:bigint;pk;id;notnull"`
	FileName string `gondolier:"type:character varying(255);notnull"`
}

type testPost struct {
	Id      uint64 `gondolier:"type:bigint;pk;id;notnull"`
	Post    string `gondolier:"type:character varying(255);notnull"`
	User    uint64 `gondolier:"type:bigint;fk:testUser;notnull"`
	Picture uint64 `gondolier:"type:bigint;fk:testPicture;null"`
}

func TestPostgresCreateTable(t *testing.T) {
	testCleanDb()
	postgres := NewPostgres("public")
	Use(testdb, postgres)
	Model(testUser{}, testPost{}, testPicture{})
	Migrate()

	if !postgres.tableExists("test_post") {
		t.Fatal("Table must have been created: test_post")
	}

	if !postgres.tableExists("test_user") {
		t.Fatal("Table must have been created: test_user")
	}

	if !postgres.tableExists("test_picture") {
		t.Fatal("Table must have been created: test_picture")
	}
}

func TestPostgresDropTable(t *testing.T) {
	testCleanDb()

	if _, err := testdb.Exec(`CREATE TABLE "test_user" ("id" bigint not null)`); err != nil {
		t.Fatal(err)
	}

	postgres := NewPostgres("public")
	Use(testdb, postgres)
	Drop(testUser{})

	if postgres.tableExists("test_user") {
		t.Fatal("Table must have been dropped")
	}
}

func TestPostgresDropTableNotExists(t *testing.T) {
	testCleanDb()
	postgres := NewPostgres("public")
	Use(testdb, postgres)
	Drop(testUser{})

	if postgres.tableExists("test_user") {
		t.Fatal("Table must have been dropped")
	}
}

func testCleanDb() {
	testdb.Exec(`DROP SEQUENCE IF EXISTS "test_post_id_seq"`)
	testdb.Exec(`DROP SEQUENCE IF EXISTS "test_user_id_seq"`)
	testdb.Exec(`DROP SEQUENCE IF EXISTS "test_picture_id_seq"`)
	testdb.Exec(`DROP TABLE IF EXISTS "test_post"`)
	testdb.Exec(`DROP TABLE IF EXISTS "test_user"`)
	testdb.Exec(`DROP TABLE IF EXISTS "test_picture"`)
}
