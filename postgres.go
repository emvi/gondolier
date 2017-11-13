package gondolier

import (
	"strings"
)

type Postgres struct {
	Schema      string
	DropColumns bool

	createSeq []string
	alterSeq  []string
	createFK  []string
}

func (m *Postgres) Migrate(metaModels []MetaModel) {
	// create or update table
	for _, model := range metaModels {
		m.migrate(&model)
	}

	// create foreign keys
	for _, fk := range m.createFK {
		if _, err := db.Exec(fk); err != nil {
			panic(err)
		}
	}

	// drop foreign keys
	// TODO

	// reset
	m.createFK = make([]string, 0)
}

func (m *Postgres) DropTable(name string) {
	name = naming.Get(name)

	if _, err := db.Exec(`DROP TABLE IF EXISTS "` + name + `"`); err != nil {
		panic(err)
	}
}

func (m *Postgres) migrate(model *MetaModel) {
	if !m.tableExists(model.ModelName) {
		m.createTable(model)
	} else {
		m.updateTable(model)
		m.dropColumns(model)
	}
}

func (m *Postgres) tableExists(name string) bool {
	name = naming.Get(name)

	rows, err := db.Query(`SELECT EXISTS (SELECT 1
	   FROM information_schema.tables
	   WHERE table_schema = $1
	   AND table_name = $2)`, m.Schema, name)

	if err != nil {
		panic(err)
	}

	var exists bool
	rows.Next()

	if err := rows.Scan(&exists); err != nil {
		panic(err)
	}

	return exists
}

func (m *Postgres) columnExists(name string) bool {
	// TODO
	return false
}

func (m *Postgres) sequenceExists(name string) bool {
	name = naming.Get(name)

	rows, err := db.Query(`SELECT EXISTS (SELECT 1
	   FROM pg_class
	   WHERE relkind = 'S'
	   AND oid::regclass::text = quote_ident($1))`, name)

	if err != nil {
		panic(err)
	}

	var exists bool
	rows.Next()

	if err := rows.Scan(&exists); err != nil {
		panic(err)
	}

	return exists
}

func (m *Postgres) foreignKeyExists(tableName, fkName string) bool {
	tableName = naming.Get(tableName)
	fkName = naming.Get(fkName)

	rows, err := db.Query(`SELECT EXISTS (SELECT 1
		FROM information_schema.table_constraints
		WHERE constraint_name = $1
		AND table_name = $2)`, fkName, tableName)

	if err != nil {
		panic(err)
	}

	var exists bool
	rows.Next()

	if err := rows.Scan(&exists); err != nil {
		panic(err)
	}

	return exists
}

func (m *Postgres) createTable(model *MetaModel) {
	name := naming.Get(model.ModelName)
	sql := `CREATE TABLE IF NOT EXISTS "` + name + `" (` + m.getColumns(model) + `)`

	// create sequences if required
	for _, seq := range m.createSeq {
		if _, err := db.Exec(seq); err != nil {
			panic(err)
		}
	}

	// create table
	if _, err := db.Exec(sql); err != nil {
		panic(err)
	}

	// alter sequence if required
	for _, seq := range m.alterSeq {
		if _, err := db.Exec(seq); err != nil {
			panic(err)
		}
	}

	// reset
	m.createSeq = make([]string, 0)
	m.alterSeq = make([]string, 0)
}

func (m *Postgres) updateTable(model *MetaModel) {
	// TODO
}

func (m *Postgres) dropColumns(model *MetaModel) {
	// TODO
}

func (m *Postgres) getColumns(model *MetaModel) string {
	columns := ""

	for _, field := range model.Fields {
		columns += `"` + naming.Get(field.Name) + `" ` + m.getTags(model.ModelName, &field) + `,`
	}

	return columns[:len(columns)-1]
}

func (m *Postgres) getTags(modelName string, field *MetaField) string {
	tags := make([]string, 5)

	for _, tag := range field.Tags {
		key := strings.ToLower(tag.Name)
		value := strings.ToLower(tag.Value)

		if key == "type" {
			tags[0] = tag.Value
		} else if key == "default" {
			tags[1] = "DEFAULT "

			if value == "nextval(seq)" {
				tags[1] += "nextval('" + m.getSequenceName(modelName, field.Name) + "'::regclass)"
			} else {
				tags[1] += value
			}
		} else if value == "notnull" || value == "not null" {
			tags[2] = "NOT NULL"
		} else if value == "null" {
			tags[2] = "NULL"
		} else if key == "seq" {
			m.addSequence(modelName, field.Name, value)
		} else if value == "id" {
			// id is a shortcut for seq + default
			m.addSequence(modelName, field.Name, "1,1,-,-,1")
			tags[1] = "DEFAULT nextval('" + m.getSequenceName(modelName, field.Name) + "'::regclass)"
		} else if value == "pk" || value == "primary key" {
			tags[3] = "PRIMARY KEY"
		} else if value == "unique" {
			tags[4] = "UNIQUE"
		} else if key == "fk" || key == "foreign key" {
			// value must be case sensitive here
			m.addForeignKey(modelName, field.Name, tag.Value)
		} else {
			name := ""

			if key == "" {
				name = value
			} else {
				name = key + ":" + value
			}

			panic("Unknown tag '" + name + "' for model '" + modelName + "'")
		}
	}

	return strings.Join(tags, " ")
}

func (m *Postgres) addSequence(modelName, columnName, info string) {
	// create sequence
	infos := strings.Split(info, ",")

	if len(infos) != 5 {
		panic("Five arguments must be specified for seq in model '" + modelName + "': start, increment, min, max, cache")
	}

	name := m.getSequenceName(modelName, columnName)
	seq := `CREATE SEQUENCE IF NOT EXISTS "` + name + `"
		START WITH ` + infos[0] + `
		INCREMENT BY ` + infos[1]

	if infos[2] == "-" {
		seq += " NO MINVALUE"
	} else {
		seq += " MINVALUE " + infos[2]
	}

	if infos[3] == "-" {
		seq += " NO MAXVALUE"
	} else {
		seq += " MAXVALUE " + infos[3]
	}

	if infos[4] != "-" {
		seq += " CACHE " + infos[4]
	}

	m.createSeq = append(m.createSeq, seq)

	// create owned by table
	modelName = naming.Get(modelName)
	columnName = naming.Get(columnName)
	alterSeq := `ALTER SEQUENCE "` + name + `"
		OWNED BY "` + modelName + `"."` + columnName + `"`
	m.alterSeq = append(m.alterSeq, alterSeq)
}

func (m *Postgres) getSequenceName(modelName, columnName string) string {
	modelName = naming.Get(modelName)
	columnName = naming.Get(columnName)
	return modelName + "_" + columnName + "_seq"
}

func (m *Postgres) addForeignKey(modelName, columnName, info string) {
	infos := strings.Split(info, ".")

	if len(infos) != 2 {
		panic("Two arguments must be specified for fk in model '" + modelName + "': ReferencedModel.ReferencedAttribute")
	}

	tableName := naming.Get(modelName)
	columnName = naming.Get(columnName)
	refTableName := naming.Get(infos[0])
	refColumnName := naming.Get(infos[1])
	fkName := m.getForeignKeyName(modelName, infos[0])
	alterFk := `ALTER TABLE "` + tableName + `"
		ADD CONSTRAINT "` + fkName + `"
		FOREIGN KEY ("` + columnName + `")
		REFERENCES "` + refTableName + `"("` + refColumnName + `")`
	m.createFK = append(m.createFK, alterFk)
}

func (m *Postgres) getForeignKeyName(modelName, refObjName string) string {
	modelName = naming.Get(modelName)
	refObjName = naming.Get(refObjName)
	return modelName + "_" + refObjName + "_fk"
}
