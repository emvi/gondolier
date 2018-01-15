package gondolier

import (
	"database/sql"
	"strings"
)

var (
	db         *sql.DB
	migrator   Migrator
	naming     = NameSchema(&SnakeCase{})
	metaModels = make([]MetaModel, 0)
)

// Migrator interface used to migrate a database schema for a specific database.
type Migrator interface {
	Migrate([]MetaModel)
	DropTable(string)
}

// NameSchema interface used to translate model names to schema names.
type NameSchema interface {
	Get(string) string
}

// Use sets the database connection and migrator.
func Use(conn *sql.DB, m Migrator) {
	db = conn
	migrator = m
}

// Naming sets the naming pattern used for migration. Default is snake case.
//
// Example:
//  Naming(SnakeCase)
func Naming(schema NameSchema) {
	if schema == nil {
		panic("Name schema must not be nil")
	}

	naming = schema
}

// Model adds one or more objects for migration.
// The objects can be passed as references, values or mixed.
// This function might panic if an invalid model is used.
//
// Example:
//  Model(&MyModel{}, AnotherModel{})
func Model(models ...interface{}) {
	for _, model := range models {
		if !modelExists(model) {
			metaModels = append(metaModels, buildMetaModel(model))
		}
	}
}

// Migrate migrates models added previously using Model().
// The database connection and migrator must be set before by calling Use().
//
// Example:
//  Use(Postgres)
//  Model(MyModel{}, AnotherModel{})
//  Migrate()
func Migrate() {
	checkSetup()
	migrator.Migrate(metaModels)
	reset()
}

// Drop drops tables for given objects if they exist.
// The database connection and migrator must be set before by calling Use().
// The objects can be passed as references, values or mixed.
// This function might panic if an invalid model is used or the tables cannot be dropped.
//
// Example:
//  Drop(&MyModel{}, AnotherModel{})
func Drop(models ...interface{}) {
	checkSetup()

	for _, model := range models {
		metaModel := buildMetaModel(model)
		migrator.DropTable(metaModel.ModelName)
	}
}

func modelExists(model interface{}) bool {
	name := strings.ToLower(getModelName(model))

	for _, metaModel := range metaModels {
		if name == strings.ToLower(metaModel.ModelName) {
			return true
		}
	}

	return false
}

func checkSetup() {
	if db == nil {
		panic("No database connection was set, call Use(connection, migrator) to set one")
	}

	if migrator == nil {
		panic("No migrator was set, call Use(connection, migrator) to select one")
	}

	if naming == nil {
		panic("No naming was set, call Naming(naming) to set one")
	}
}

func reset() {
	metaModels = make([]MetaModel, 0)
}
