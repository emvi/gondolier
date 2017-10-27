package model

import (
	"testing"
)

type testModel struct {
	Id    uint64 `model:"type:integer;primary_key;notnull;;"`
	Name  string `model:"type:varchar;unique"`
	Age   int    `model:"type:integer;notnull"`
	Array []int  `model:"type:integer[]"`
}

type testInvalidTypesModel struct {
	Ignored          bool `model:"-"`
	IgnoredToo       bool
	Unknown          struct{ Name string } `model:"type:struct"`
	UnknownToo       *int                  `model:"type:integer"`
	UnknownInterface interface{}           `model:"type:interface"`
}

func TestBuildMetaModel(t *testing.T) {
	meta := buildMetaModel(&testModel{})

	if meta.ModelName != "testModel" {
		t.Fatal("Name must be testModel")
	}
}

func TestGetModelName(t *testing.T) {
	name := getModelName(&testModel{})

	if name != "testModel" {
		t.Fatalf("Model name must be testModel, but was %v", name)
	}

	name = getModelName(testModel{})

	if name != "testModel" {
		t.Fatalf("Model name must be testModel, but was %v", name)
	}
}

func TestGetModelNameStructOnly(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Calling getModelName with invalid type must panic")
		}
	}()

	val := 42
	getModelName(val)
	getModelName(&val)
}

func TestGetModelFields(t *testing.T) {
	fields := getModelFields(&testModel{})

	if len(fields) != 4 {
		t.Fatalf("All fields must be returned: %v", len(fields))
	}

	if fields[0].Name != "Id" ||
		fields[1].Name != "Name" ||
		fields[2].Name != "Age" ||
		fields[3].Name != "Array" {
		t.Fatal("Field names must be correct")
	}
}

func TestParseTag(t *testing.T) {
	tags := parseTag("type:varchar(20);primary_key;notnull")

	if len(tags) != 3 {
		t.Fatal("All elements must be returned")
	}

	if tags[0].Name != "type" || tags[0].Value != "varchar(20)" ||
		tags[1].Name != "" || tags[1].Value != "primary_key" ||
		tags[2].Name != "" || tags[2].Value != "notnull" {
		t.Fatal("Tag elements must be correct")
	}
}
