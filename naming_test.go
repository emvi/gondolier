package gondolier

import (
	"testing"
)

func TestToSnakeCase(t *testing.T) {
	names := [][]string{
		{"", ""},
		{"a", "a"},
		{"myModelName", "my_model_name"},
		{"APISnakeNAME", "api_snake_name"},
		{"First", "first"},
		{"SNake", "s_nake"},
		{"with_underscore", "with_underscore"},
		{"snake_Id", "snake_id"},
		{"snake_ID", "snake_id"},
		{"snakeID", "snake_id"},
		{"myLITTLEPony", "my_little_pony"},
		{"MYSnake_case", "my_snake_case"},
		{"WOOFWoof", "woof_woof"},
		{"Space to underscore", "space_to_underscore"},
	}
	namesake := SnakeCase{}

	for _, name := range names {
		got := namesake.Get(name[0])

		if got != name[1] {
			t.Fatalf("Expected snake case name %v but got %v for input %v", name[1], got, name[0])
		}
	}
}
