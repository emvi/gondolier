package model

import (
	"testing"
)

func TestToSnakeCase(t *testing.T) {
	names := [][]string{[]string{"", ""},
		[]string{"a", "a"},
		[]string{"myModelName", "my_model_name"},
		[]string{"APISnakeNAME", "api_snake_name"},
		[]string{"First", "first"},
		[]string{"SNake", "s_nake"},
		[]string{"with_underscore", "with_underscore"},
		[]string{"snake_Id", "snake_id"},
		[]string{"snake_ID", "snake_id"},
		[]string{"snakeID", "snake_id"},
		[]string{"myLITTLEPony", "my_little_pony"},
		[]string{"MYSnake_case", "my_snake_case"},
		[]string{"WOOFWoof", "woof_woof"},
		[]string{"Space to underscore", "space_to_underscore"}}

	for _, name := range names {
		got := toSnakeCase(name[0])

		if got != name[1] {
			t.Fatalf("Expected snake case name %v but got %v for input %v", name[1], got, name[0])
		}
	}
}
