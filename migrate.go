package model

import (
	"strings"
)

const (
	Postgres = "postgres"
)

var (
	database = "" // no default
)

// Sets the database type used for migration.
func Use(db string) {
	db = strings.ToLower(db)

	if db != Postgres {
		panic("The database '" + db + "' is not supported")
	}

	database = db
}

// Adds one or more models for migration. Can be passed as references to a structs or the structs directly or mixed.
// Returns true if the model was added for migration or false if it was already.
// This function might panic if an invalid model is passed.
// For example: Model(&MyModel{}, AnotherModel{})
func Model(model ...interface{}) bool {
	return true
}

// Migrates models added previously.
func Migrate() {

}
