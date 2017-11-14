package gondolier

import (
	"testing"
)

type testUser struct {
	Id      uint64 `gondolier:"type:bigint;pk;seq:1,1,-,-,1;default:nextval(seq);notnull"`
	Name    string `gondolier:"type:character varying(255);notnull;unique"`
	Age     uint   `gondolier:"type:integer;notnull"`
	Picture uint64 `gondolier:"type:bigint;fk:testPicture.Id;null"`
}

type testPicture struct {
	Id       uint64 `gondolier:"type:bigint;pk;id;notnull"`
	FileName string `gondolier:"type:character varying(255);notnull"`
}

type testPost struct {
	Id      uint64 `gondolier:"type:bigint;pk;id;notnull"`
	Post    string `gondolier:"type:character varying(255);notnull"`
	User    uint64 `gondolier:"type:bigint;fk:testUser.Id;notnull"`
	Picture uint64 `gondolier:"type:bigint;fk:testPicture.Id;null"`
}

type testArticle struct {
	Id            uint64   `json:"id" gondolier:"type:bigint;pk;id;notnull"`
	Filename      string   `json:"filename" gondolier:"type:character varying(255);notnull"`
	Tags          []string `json:"tags" gondolier:"type:character varying(255)[]"`
	Views         uint     `json:"views" gondolier:"type:integer;notnull"`
	WIP           bool     `json:"wip" gondolier:"type:boolean;notnull"`
	ReadEveryone  bool     `db:"read_everyone" json:"read_everyone" gondolier:"type:boolean;notnull"`
	WriteEveryone bool     `db:"write_everyone" json:"write_everyone" gondolier:"type:boolean;notnull"`

	SomeSlice []int `db:"-" json:"some_slice" gondolier:"-"`
}

func TestPostgresCreateTable(t *testing.T) {
	testCleanDb()
	postgres := &Postgres{Schema: "public"}
	Use(testdb, postgres)
	Model(testUser{}, testPost{}, testPicture{}, testArticle{})
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

	if !postgres.tableExists("test_article") {
		t.Fatal("Table must have been created: test_article")
	}

	if !postgres.sequenceExists("test_post_id_seq") {
		t.Fatal("Sequence must have been created: test_post_id_seq")
	}

	if !postgres.sequenceExists("test_user_id_seq") {
		t.Fatal("Sequence must have been created: test_user_id_seq")
	}

	if !postgres.sequenceExists("test_picture_id_seq") {
		t.Fatal("Sequence must have been created: test_picture_id_seq")
	}

	if !postgres.sequenceExists("test_article_id_seq") {
		t.Fatal("Sequence must have been created: test_article_id_seq")
	}

	if !postgres.foreignKeyExists("test_user", "test_user_test_picture_fk") {
		t.Fatal("Foreign key must have been created: test_user_test_picture_fk")
	}

	if !postgres.foreignKeyExists("test_post", "test_post_test_user_fk") {
		t.Fatal("Foreign key must have been created: test_post_test_user_fk")
	}

	if !postgres.foreignKeyExists("test_post", "test_post_test_picture_fk") {
		t.Fatal("Foreign key must have been created: test_post_test_picture_fk")
	}
}

func TestPostgresDropTable(t *testing.T) {
	testCleanDb()

	if _, err := testdb.Exec(`CREATE TABLE "test_user" ("id" bigint not null)`); err != nil {
		t.Fatal(err)
	}

	postgres := &Postgres{Schema: "public"}
	Use(testdb, postgres)
	Drop(testUser{})

	if postgres.tableExists("test_user") {
		t.Fatal("Table must have been dropped")
	}
}

func TestPostgresDropTableNotExists(t *testing.T) {
	testCleanDb()
	postgres := &Postgres{Schema: "public"}
	Use(testdb, postgres)
	Drop(testUser{})

	if postgres.tableExists("test_user") {
		t.Fatal("Table must have been dropped")
	}
}

func TestPostgresDropColumn(t *testing.T) {
	testCleanDb()

	if _, err := testdb.Exec(`CREATE TABLE "test_user"
		("id" bigint not null, "drop_me" text not null)`); err != nil {
		t.Fatal(err)
	}

	postgres := &Postgres{Schema: "public", DropColumns: true}
	Use(testdb, postgres)
	Model(testUser{})
	Migrate()

	if postgres.columnExists("test_user", "drop_me") {
		t.Fatal("Column 'drop_me' should not exist anymore")
	}

	if !postgres.columnExists("test_user", "id") {
		t.Fatal("Column 'id' must still exist")
	}
}

func testCleanDb() {
	testdb.Exec(`DROP SEQUENCE IF EXISTS "test_post_id_seq"`)
	testdb.Exec(`DROP SEQUENCE IF EXISTS "test_user_id_seq"`)
	testdb.Exec(`DROP SEQUENCE IF EXISTS "test_picture_id_seq"`)
	testdb.Exec(`DROP SEQUENCE IF EXISTS "test_article_id_seq"`)
	testdb.Exec(`DROP TABLE IF EXISTS "test_post"`)
	testdb.Exec(`DROP TABLE IF EXISTS "test_user"`)
	testdb.Exec(`DROP TABLE IF EXISTS "test_picture"`)
	testdb.Exec(`DROP TABLE IF EXISTS "test_article"`)
}
