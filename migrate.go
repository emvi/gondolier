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

type Migrator interface {
	Migrate(*sql.Tx, []MetaModel)
	DropTable(*sql.Tx, string)
}

type NameSchema interface {
	Get(string) string
}

// Sets the database connection and migrator used for migration.
func Use(conn *sql.DB, m Migrator) {
	db = conn
	migrator = m
}

// Sets the naming used for migration. Default is snake case.
// Example: Naming(SnakeCase)
func Naming(schema NameSchema) {
	if schema == nil {
		panic("Name schema must not be nil")
	}

	naming = schema
}

// Adds one or more models for migration.
// Can be passed as references to a structs or the structs directly or mixed.
// This function might panic if an invalid model is passed.
// Example: Model(&MyModel{}, AnotherModel{})
func Model(models ...interface{}) {
	for _, model := range models {
		if !modelExists(model) {
			metaModels = append(metaModels, buildMetaModel(model))
		}
	}
}

// Migrates models added previously using Model().
// The database connection and migrator must be set.
// Example:
//
// Use(Postgres)
// Model(MyModel{}, AnotherModel{})
// Migrate()
func Migrate() {
	checkSetup()
	tx := begin()
	defer rollback(tx)
	migrator.Migrate(tx, metaModels)
	reset()
	commit(tx)
}

// Drops tables for given models if exists.
// The database connection and migrator must be set.
// Can be passed as references to a structs or the structs directly or mixed.
// This function might panic if an invalid model is passed or the tables cannot be dropped.
// Example: Drop(&MyModel{}, AnotherModel{})
func Drop(models ...interface{}) {
	checkSetup()
	tx := begin()
	defer rollback(tx)

	for _, model := range models {
		metaModel := buildMetaModel(model)
		migrator.DropTable(tx, metaModel.ModelName)
	}

	commit(tx)
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

func begin() *sql.Tx {
	tx, err := db.Begin()

	if err != nil {
		panic(err)
	}

	return tx
}

func rollback(tx *sql.Tx) {
	if r := recover(); r != nil {
		if err := tx.Rollback(); err != nil {
			panic(err)
		}
	}
}

func commit(tx *sql.Tx) {
	if err := tx.Commit(); err != nil {
		panic(err)
	}
}

func reset() {
	metaModels = make([]MetaModel, 0)
}
