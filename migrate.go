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

// Interface used to migrate a database schema from model.
type Migrator interface {
	Migrate([]MetaModel)
	DropTable(string)
}

// Interface used to translate model names to schema names.
type NameSchema interface {
	Get(string) string
}

// Use sets the database connection and migrator used for migration.
func Use(conn *sql.DB, m Migrator) {
	db = conn
	migrator = m
}

// Naming sets the naming used for migration. Default is snake case.
//
// Example:
//  Naming(SnakeCase)
func Naming(schema NameSchema) {
	if schema == nil {
		panic("Name schema must not be nil")
	}

	naming = schema
}

// Model adds one or more models for migration.
// Can be passed as references to a structs or the structs directly or mixed.
// This function might panic if an invalid model is passed.
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
// The database connection and migrator must be set.
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

// Drop drops tables for given models if exists.
// The database connection and migrator must be set.
// Can be passed as references to a structs or the structs directly or mixed.
// This function might panic if an invalid model is passed or the tables cannot be dropped.
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
