package gondolier

import (
	"database/sql"
	"testing"
)

type testModelA struct {
	Id   uint64 `gondolier:"type:bigint;primarykey;notnull"`
	Name string `gondolier:"type:character varying(100)"`
	B    uint64 `gondolier:"type:bigint;fk(testModelB)"`
}

type testModelB struct {
	Id  uint64 `gondolier:"type:bigint;primarykey;notnull"`
	Age uint   `gondolier:"type:integer"`
}

type dummyMigrator struct {
	models []MetaModel
	drop   []string
}

func (m *dummyMigrator) Migrate(tx *sql.Tx, metaModels []MetaModel) {
	m.models = metaModels
}

func (m *dummyMigrator) DropTable(tx *sql.Tx, name string) {
	m.drop = append(m.drop, name)
}

type dummyCase struct{}

func (n *dummyCase) Get(name string) string {
	return "works"
}

func TestUse(t *testing.T) {
	if migrator != nil {
		t.Fatal("No migrator must be selected")
	}

	Use(testdb, &Postgres{})

	if migrator == nil {
		t.Fatal("Postgres must be selected")
	}
}

func testNaming(t *testing.T) {
	Naming(&dummyCase{})

	if naming.Get("") != "works" {
		t.Fatal("Name schema must have been set")
	}
}

func TestNamingNotNil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Naming must panic if nil was passed")
		}
	}()

	Naming(nil)
}

func TestModel(t *testing.T) {
	Model(&testModelA{}, testModelB{}, &testModelB{})

	if len(metaModels) != 2 {
		t.Fatal("Two models must have been added")
	}
}

func TestMigrate(t *testing.T) {
	dummy := &dummyMigrator{}
	Use(testdb, dummy)
	Model(testModelA{}, testModelB{})
	Migrate()

	if len(dummy.models) != 2 {
		t.Fatal("Translate must have been called")
	}
}

func TestMigrateNoMigrator(t *testing.T) {
	migrator = nil

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Migrate must panic if no migrator was selected")
		}
	}()

	Migrate()
}

func TestDrop(t *testing.T) {
	dummy := &dummyMigrator{}
	Use(testdb, dummy)
	Drop(testModelA{}, testModelB{})

	if len(dummy.drop) != 2 {
		t.Fatal("Two models must have been dropped")
	}

	if dummy.drop[0] != "testModelA" || dummy.drop[1] != "testModelB" {
		t.Fatalf("Two tables must have been dropped, but was %v %v", dummy.drop[0], dummy.drop[1])
	}
}
