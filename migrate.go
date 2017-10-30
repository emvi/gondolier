package gondolier

import (
	"strings"
)

var (
	migrator   Migrator // no default
	metaModels = make([]MetaModel, 0)
)

type Migrator interface {
	Migrate([]MetaModel)
	DropTable(string)
}

// Sets the database migrator used for migration.
func Use(m Migrator) {
	migrator = m
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

// Migrates models added previously using Model(). The migrator must be set before.
// Example:
//
// Use(Postgres)
// Model(MyModel{}, AnotherModel{})
// Migrate()
func Migrate() {
	migratorSet()
	migrator.Migrate(metaModels)
	reset()
}

// Drops tables for given models if exists. The migrator must be set before.
// Can be passed as references to a structs or the structs directly or mixed.
// This function might panic if an invalid model is passed.
// Example: Drop(&MyModel{}, AnotherModel{})
func Drop(models ...interface{}) {
	migratorSet()

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

func migratorSet() {
	if migrator == nil {
		panic("No migrator was set, call Use(migrator) to select one")
	}
}

func reset() {
	metaModels = make([]MetaModel, 0)
}
