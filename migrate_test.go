package model

import (
	"testing"
)

type testModelA struct {
	Id   uint64 `model:"type:bigint;primarykey;notnull"`
	Name string `model:"type:character varying(100)"`
	B    uint64 `model:"type:bigint;fk(testModelB)"`
}

type testModelB struct {
	Id  uint64 `model:"type:bigint;primarykey;notnull"`
	Age uint   `model:"type:integer"`
}

type dummyTranslator struct {
	models []MetaModel
}

func (t *dummyTranslator) Translate(metaModels []MetaModel) string {
	t.models = metaModels
	return ""
}

func TestUse(t *testing.T) {
	if translator != nil {
		t.Fatal("No translator must be selected")
	}

	Use("Postgres")

	if translator == nil {
		t.Fatal("Postgres must be selected")
	}
}

func TestUseUnknownDb(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Fatal("Calling Use with unknown database must panic")
		}
	}()

	Use("unknown")
}

func TestModel(t *testing.T) {
	Model(&testModelA{}, testModelB{}, &testModelB{})

	if len(metaModels) != 2 {
		t.Fatal("Two models must have been added")
	}
}

func TestMigrate(t *testing.T) {
	dummy := &dummyTranslator{}
	translator = dummy
	Model(testModelA{}, testModelB{})
	Migrate()

	if len(dummy.models) != 2 {
		t.Fatal("Translate must have been called")
	}
}

func TestMigrateNoTranslator(t *testing.T) {
	translator = nil

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Migrate must panic if no translator was selected")
		}
	}()

	Migrate()
}
