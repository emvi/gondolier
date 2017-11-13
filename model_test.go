package gondolier

import (
	"testing"
	"time"
)

type testModel struct {
	Id    uint64    `gondolier:"type:integer;primarykey;notnull;;"` // accept multiple ;
	Name  string    `gondolier:"type:varchar;unique"`
	Age   int       `gondolier:"type:integer;notnull"`
	Array []int     `gondolier:"type:integer[]"`
	Date  time.Time `gondolier:"type:timestamp"`
}

type testInvalidTypesModel struct {
	Ignored          bool `gondolier:"-"`
	IgnoredToo       bool
	Unknown          struct{ Name string } `gondolier:"type:struct"`
	UnknownToo       *int                  `gondolier:"type:integer"`
	UnknownInterface interface{}           `gondolier:"type:interface"`
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

	if len(fields) != 5 {
		t.Fatalf("All fields must be returned: %v", len(fields))
	}

	if fields[0].Name != "Id" ||
		fields[1].Name != "Name" ||
		fields[2].Name != "Age" ||
		fields[3].Name != "Array" ||
		fields[4].Name != "Date" {
		t.Fatal("Field names must be correct")
	}

	fields = getModelFields(testModel{})

	if len(fields) != 5 {
		t.Fatalf("All fields must be returned: %v", len(fields))
	}
}

func TestParseTag(t *testing.T) {
	tags := parseTag("type:varchar(20);primarykey;notnull")

	if len(tags) != 3 {
		t.Fatal("All elements must be returned")
	}

	if tags[0].Name != "type" || tags[0].Value != "varchar(20)" ||
		tags[1].Name != "" || tags[1].Value != "primarykey" ||
		tags[2].Name != "" || tags[2].Value != "notnull" {
		t.Fatal("Tag elements must be correct")
	}
}
