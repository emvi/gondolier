package model

import (
	"testing"
)

type testModel struct {
	Id   uint64 `model:"primary_key"`
	Name string `model:"type:varchar(20);unique"`
	Age  int    `model:"type:integer`
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

func TestGetModelFields(t *testing.T) {
	fields := getModelFields(&testModel{})

	if len(fields) != 3 {
		t.Fatal("All fields must be returned")
	}

	if fields[0].Name != "Id" || fields[1].Name != "Name" || fields[2].Name != "Age" {
		t.Fatal("Field names must be correct")
	}
}
