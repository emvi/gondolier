package model

import (
	"testing"
)

func TestToSnakeCase(t *testing.T) {
	name := toSnakeCase("myModelName")

	if name != "my_model_name" {
		t.Fatalf("Model name must by my_model_name but was %v", name)
	}

	name = toSnakeCase("MYmodelNAME")

	if name != "my_model_name" {
		t.Fatalf("Model name must by my_model_name but was %v", name)
	}
}
