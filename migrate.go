package model

import (
	"strings"
)

const (
	Postgres = "postgres"
)

var (
	translator Translator // no default
	metaModels = make([]MetaModel, 0)
)

// A translator translates the meta model into executable SQL statements in right order.
type Translator interface {
	// Translate takes the meta model and translates it to executable SQL statements,
	// which are returned as one string.
	Translate([]MetaModel) string
}

// Sets the database type used for migration.
func Use(db string) {
	db = strings.ToLower(db)

	if db == Postgres {
		translator = &postgresTranslator{}
	} else {
		panic("The database '" + db + "' is not supported")
	}
}

// Adds one or more models for migration. Can be passed as references to a structs or the structs directly or mixed.
// This function might panic if an invalid model is passed.
// For example: Model(&MyModel{}, AnotherModel{})
func Model(models ...interface{}) {
	for _, model := range models {
		if !modelExists(model) {
			metaModels = append(metaModels, buildMetaModel(model))
		}
	}
}

// Migrates models added previously.
func Migrate() {
	if translator == nil {
		panic("No translator was selected, call Use(database) to select a translator")
	}

	translator.Translate(metaModels)
	reset()
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

func reset() {
	metaModels = make([]MetaModel, 0)
}
