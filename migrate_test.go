package gondolier

import (
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

func (m *dummyMigrator) Migrate(metaModels []MetaModel) {
	m.models = metaModels
}

func (m *dummyMigrator) DropTable(name string) {
	m.drop = append(m.drop, name)
}

func TestUse(t *testing.T) {
	if migrator != nil {
		t.Fatal("No migrator must be selected")
	}

	Use(&Postgres{})

	if migrator == nil {
		t.Fatal("Postgres must be selected")
	}
}

func TestModel(t *testing.T) {
	Model(&testModelA{}, testModelB{}, &testModelB{})

	if len(metaModels) != 2 {
		t.Fatal("Two models must have been added")
	}
}

func TestMigrate(t *testing.T) {
	dummy := &dummyMigrator{}
	Use(dummy)
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
	Use(dummy)
	Drop(testModelA{}, testModelB{})

	if len(dummy.drop) != 2 {
		t.Fatal("Two models must have been dropped")
	}

	if dummy.drop[0] != "testModelA" || dummy.drop[1] != "testModelB" {
		t.Fatalf("Two tables must have been dropped, but was %v %v", dummy.drop[0], dummy.drop[1])
	}
}
