package gondolier

import (
	"testing"
)

type testUser struct {
	Id      uint64 `gondolier:"type:bigint;pk;seq:1,1,-,-,1;default:nextval(seq);notnull"`
	Name    string `gondolier:"type:varchar(255);notnull;unique"`
	Age     uint   `gondolier:"type:integer;notnull"`
	Picture uint64 `gondolier:"type:bigint;fk:test_picture.id;null"`
}

type testPicture struct {
	Id       uint64 `gondolier:"type:bigint;id"`
	FileName string `gondolier:"type:varchar(255);notnull"`
}

type testPost struct {
	Id      uint64 `gondolier:"type:bigint;id"`
	Post    string `gondolier:"type:varchar(255);notnull"`
	User    uint64 `gondolier:"type:bigint;fk:test_user.id;notnull"`
	Picture uint64 `gondolier:"type:bigint;fk:test_picture.id;null"`
}

type testArticle struct {
	Id            uint64   `gondolier:"type:bigint;id"`
	Filename      string   `gondolier:"type:varchar(255);notnull"`
	Tags          []string `gondolier:"type:varchar(255)[]"`
	Views         uint     `gondolier:"type:integer;notnull"`
	WIP           bool     `gondolier:"type:boolean;notnull"`
	ReadEveryone  bool     `gondolier:"type:boolean;notnull"`
	WriteEveryone bool     `gondolier:"type:boolean;notnull"`

	SomeSlice []int `gondolier:"-"`
}

type testDropColumn struct {
	Id uint64 `gondolier:"type:bigint;id"`
}

type testAddColumn struct {
	Id        uint64 `gondolier:"type:bigint;id"`
	NewColumn string `gondolier:"type:varchar(255)"`
}

type testUpdateColumn struct {
	Column string `gondolier:"type:varchar(255);default:'default';notnull;unique;pk"`
}

type testUpdateColumnReduce struct {
	Column string `gondolier:"type:text"`
}

type testUpdateColumnSeq struct {
	Column int `gondolier:"type:integer;seq:1,1,-,-,1;default:nextval(seq)"`
}

type testUpdateColumnSeqReduce struct {
	Column int `gondolier:"type:integer"`
}

type testUpdateColumnFk struct {
	Fk uint64 `gondolier:"type:bigint;fk:test_other.id"`
}

type testUpdateColumnFkReduce struct {
	Fk uint64 `gondolier:"type:bigint"`
}

func TestPostgresCreateTable(t *testing.T) {
	testCleanDb()
	t.Log("--- TestPostgresCreateTable ---")

	postgres := &Postgres{Schema: "public", Log: true}
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

	if !postgres.foreignKeyExists("test_user", "test_user_test_picture_id_fk") {
		t.Fatal("Foreign key must have been created: test_user_test_picture_fk")
	}

	if !postgres.foreignKeyExists("test_post", "test_post_test_user_id_fk") {
		t.Fatal("Foreign key must have been created: test_post_test_user_fk")
	}

	if !postgres.foreignKeyExists("test_post", "test_post_test_picture_id_fk") {
		t.Fatal("Foreign key must have been created: test_post_test_picture_fk")
	}
}

func TestPostgresDropTable(t *testing.T) {
	testCleanDb()
	t.Log("--- TestPostgresDropTable ---")

	if _, err := testdb.Exec(`CREATE TABLE "test_user" ("id" bigint not null)`); err != nil {
		t.Fatal(err)
	}

	postgres := &Postgres{Schema: "public", Log: true}
	Use(testdb, postgres)
	Drop(testUser{})

	if postgres.tableExists("test_user") {
		t.Fatal("Table must have been dropped")
	}
}

func TestPostgresDropTableNotExists(t *testing.T) {
	testCleanDb()
	t.Log("--- TestPostgresDropTableNotExists ---")

	postgres := &Postgres{Schema: "public", Log: true}
	Use(testdb, postgres)
	Drop(testUser{})

	if postgres.tableExists("test_user") {
		t.Fatal("Table must have been dropped")
	}
}

func TestPostgresDropColumn(t *testing.T) {
	testCleanDb()
	t.Log("--- TestPostgresDropColumn ---")

	if _, err := testdb.Exec(`CREATE TABLE "test_drop_column"
		("id" bigint not null, "drop_me" text not null)`); err != nil {
		t.Fatal(err)
	}

	postgres := &Postgres{Schema: "public", DropColumns: true, Log: true}
	Use(testdb, postgres)
	Model(testDropColumn{})
	Migrate()

	if postgres.columnExists("test_drop_column", "drop_me") {
		t.Fatal("Column 'drop_me' should not exist anymore")
	}

	if !postgres.columnExists("test_drop_column", "id") {
		t.Fatal("Column 'id' must still exist")
	}
}

func TestPostgresAddColumn(t *testing.T) {
	testCleanDb()
	t.Log("--- TestPostgresAddColumn ---")

	if _, err := testdb.Exec(`CREATE TABLE "test_add_column" ("id" bigint not null)`); err != nil {
		t.Fatal(err)
	}

	postgres := &Postgres{Schema: "public", Log: true}
	Use(testdb, postgres)
	Model(testAddColumn{})
	Migrate()

	if !postgres.columnExists("test_add_column", "new_column") {
		t.Fatal("Column 'new_column' must exist")
	}
}

func TestPostgresUpdateColumn(t *testing.T) {
	testCleanDb()
	t.Log("--- TestPostgresUpdateColumn ---")

	if _, err := testdb.Exec(`CREATE TABLE "test_update_column"
		("column" text)`); err != nil {
		t.Fatal(err)
	}

	postgres := &Postgres{Schema: "public", Log: true}
	Use(testdb, postgres)
	Model(testUpdateColumn{})
	Migrate()
	istype := postgres.getColumnType("test_update_column", "column")

	if istype != "character varying" {
		t.Fatalf("Type must be character varying, but was %v", istype)
	}

	if postgres.isNullable("test_update_column", "column") {
		t.Fatal("Column must not be nullable")
	}

	if !postgres.constraintExists("test_update_column_pkey") {
		t.Fatal("Primary key constraint must exist")
	}

	if !postgres.constraintExists("test_update_column_column_key") {
		t.Fatal("Unique constraint must exist")
	}
}

func TestPostgresUpdateColumnReduce(t *testing.T) {
	testCleanDb()
	t.Log("--- TestPostgresUpdateColumnReduce ---")

	if _, err := testdb.Exec(`CREATE TABLE "test_update_column_reduce"
		("column" text NOT NULL PRIMARY KEY UNIQUE)`); err != nil {
		t.Fatal(err)
	}

	if _, err := testdb.Exec(`ALTER TABLE "test_update_column_reduce"
				RENAME CONSTRAINT "test_update_column_reduce_pkey"
				TO "test_update_column_reduce_column_pkey"`); err != nil {
		t.Fatal(err)
	}

	postgres := &Postgres{Schema: "public", Log: true}
	Use(testdb, postgres)
	Model(testUpdateColumnReduce{})
	Migrate()

	if !postgres.isNullable("test_update_column_reduce", "column") {
		t.Fatal("Column must be nullable")
	}

	if postgres.constraintExists("test_update_column_reduce_pkey") {
		t.Fatal("Primary key constraint must not exist")
	}

	if postgres.constraintExists("test_update_column_column_reduce_unique") {
		t.Fatal("Unique constraint must not exist")
	}
}

func TestPostgresUpdateColumnSeq(t *testing.T) {
	testCleanDb()
	t.Log("--- TestPostgresUpdateColumnSeq ---")

	if _, err := testdb.Exec(`CREATE TABLE "test_update_column_seq"
		("column" integer)`); err != nil {
		t.Fatal(err)
	}

	postgres := &Postgres{Schema: "public", Log: true}
	Use(testdb, postgres)
	Model(testUpdateColumnSeq{})
	Migrate()

	if !postgres.sequenceExists("test_update_column_seq_column_seq") {
		t.Fatal("Sequence must exist")
	}
}

func TestPostgresUpdateColumnSeqReduce(t *testing.T) {
	testCleanDb()
	t.Log("--- TestPostgresUpdateColumnSeqReduce ---")

	if _, err := testdb.Exec(`CREATE SEQUENCE "test_update_column_seq_reduce_column_seq"
		START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1`); err != nil {
		t.Fatal(err)
	}

	if _, err := testdb.Exec(`CREATE TABLE "test_update_column_seq_reduce"
		("column" integer DEFAULT nextval('test_update_column_seq_reduce_column_seq'::regclass))`); err != nil {
		t.Fatal(err)
	}

	if _, err := testdb.Exec(`ALTER SEQUENCE "test_update_column_seq_reduce_column_seq"
		OWNED BY "test_update_column_seq_reduce"."column"`); err != nil {
		t.Fatal(err)
	}

	postgres := &Postgres{Schema: "public", Log: true}
	Use(testdb, postgres)
	Model(testUpdateColumnSeqReduce{})
	Migrate()

	if postgres.sequenceExists("test_update_column_seq_reduce_column_seq") {
		t.Fatal("Sequence must not exist")
	}
}

func TestPostgresUpdateColumnFk(t *testing.T) {
	testCleanDb()
	t.Log("--- TestPostgresUpdateColumnFk ---")

	if _, err := testdb.Exec(`CREATE TABLE "test_other"
		("id" bigint unique)`); err != nil {
		t.Fatal(err)
	}

	if _, err := testdb.Exec(`CREATE TABLE "test_update_column_fk"
		("fk" bigint)`); err != nil {
		t.Fatal(err)
	}

	postgres := &Postgres{Schema: "public", Log: true}
	Use(testdb, postgres)
	Model(testUpdateColumnFk{})
	Migrate()

	if !postgres.constraintExists("test_update_column_fk_test_other_id_fk") {
		t.Fatal("Foreign key must exist")
	}
}

func TestPostgresUpdateColumnFkReduce(t *testing.T) {
	testCleanDb()
	t.Log("--- TestPostgresUpdateColumnFkReduce ---")

	if _, err := testdb.Exec(`CREATE TABLE "test_other"
		("id" bigint unique)`); err != nil {
		t.Fatal(err)
	}

	if _, err := testdb.Exec(`CREATE TABLE "test_update_column_fk_reduce"
		("fk" bigint REFERENCES "test_other"("id"))`); err != nil {
		t.Fatal(err)
	}

	postgres := &Postgres{Schema: "public", Log: true}
	Use(testdb, postgres)
	Model(testUpdateColumnFkReduce{})
	Migrate()

	if postgres.constraintExists("test_update_column_fk_reduce_test_other_id_fk") {
		t.Fatal("Foreign key must not exist")
	}
}

func TestPostgresDontDropPk(t *testing.T) {
	testCleanDb()
	t.Log("--- TestPostgresDontDropPk ---")

	postgres := &Postgres{Schema: "public", Log: true}
	Use(testdb, postgres)
	Model(testUser{}, testPost{}, testPicture{}, testArticle{})
	Migrate()

	t.Log("------------------------------")

	Model(testUser{}, testPost{}, testPicture{}, testArticle{})
	Migrate()
}

func testCleanDb() {
	testdb.Exec(`DROP TABLE IF EXISTS "test_post"`)
	testdb.Exec(`DROP TABLE IF EXISTS "test_user"`)
	testdb.Exec(`DROP TABLE IF EXISTS "test_picture"`)
	testdb.Exec(`DROP TABLE IF EXISTS "test_article"`)
	testdb.Exec(`DROP TABLE IF EXISTS "test_drop_column"`)
	testdb.Exec(`DROP TABLE IF EXISTS "test_add_column"`)
	testdb.Exec(`DROP TABLE IF EXISTS "test_update_column"`)
	testdb.Exec(`DROP TABLE IF EXISTS "test_update_column_reduce"`)
	testdb.Exec(`DROP TABLE IF EXISTS "test_update_column_seq"`)
	testdb.Exec(`DROP TABLE IF EXISTS "test_update_column_seq_reduce"`)
	testdb.Exec(`DROP TABLE IF EXISTS "test_update_column_fk"`)
	testdb.Exec(`DROP TABLE IF EXISTS "test_update_column_fk_reduce"`)
	testdb.Exec(`DROP TABLE IF EXISTS "test_other"`)
}
