package gondolier

import (
	"database/sql"
	"testing"
	"time"
)

type testModel struct {
	Id            uint64         `gondolier:"type:integer;primarykey;notnull;;"` // accept multiple ;
	Name          string         `gondolier:"type:varchar;unique"`
	Age           int            `gondolier:"type:integer;notnull"`
	Array         []int          `gondolier:"type:integer[]"`
	Date          time.Time      `gondolier:"type:timestamp"`
	NullableField sql.NullString `gondolier:"type:text"`
}

type testModelWhitespace struct {
	Id   uint64 `  gondolier:" type:  bigint ;;; pk;;; notnull ; ; "   `
	Name string `gondolier:"	  type  :   character varying(100) ;; 	 "   `
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

	if len(fields) != 6 {
		t.Fatalf("All fields must be returned: %v", len(fields))
	}

	if fields[0].Name != "Id" ||
		fields[1].Name != "Name" ||
		fields[2].Name != "Age" ||
		fields[3].Name != "Array" ||
		fields[4].Name != "Date" ||
		fields[5].Name != "NullableField" {
		t.Fatal("Field names must be correct")
	}

	fields = getModelFields(testModel{})

	if len(fields) != 6 {
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

func TestModelWhitespace(t *testing.T) {
	meta := buildMetaModel(testModelWhitespace{})

	if meta.ModelName != "testModelWhitespace" {
		t.Fatal("Name must be testModelWhitespace")
	}

	if len(meta.Fields) != 2 {
		t.Fatalf("Model must have two fields, but was: %v", len(meta.Fields))
	}

	fields := meta.Fields

	if fields[0].Name != "Id" || fields[1].Name != "Name" {
		t.Fatal("Model fields must have proper names")
	}

	if len(fields[0].Tags) != 3 || len(fields[1].Tags) != 1 {
		t.Fatalf("Model fields must have proper tags: %v %v", len(fields[0].Tags), len(fields[1].Tags))
	}

	if fields[0].Tags[0].Name != "type" || fields[0].Tags[0].Value != "bigint" {
		t.Fatalf("First field must have type bigint: %v %v", fields[0].Tags[0].Name, fields[0].Tags[0].Value)
	}

	if fields[0].Tags[1].Name != "" || fields[0].Tags[1].Value != "pk" {
		t.Fatal("Second field must have value pk")
	}

	if fields[0].Tags[2].Name != "" || fields[0].Tags[2].Value != "notnull" {
		t.Fatal("Third field must have value notnull")
	}

	if fields[1].Tags[0].Name != "type" || fields[1].Tags[0].Value != "character varying(100)" {
		t.Fatalf("First field must have type character varying(100): %v %v", fields[1].Tags[0].Name, fields[1].Tags[0].Value)
	}
}
