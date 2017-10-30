package gondolier

import (
	"database/sql"
	"strings"
)

type Postgres struct {
	Schema    string
	sequences []string
	alterSeq  []string
}

func NewPostgres(schema string) *Postgres {
	return &Postgres{Schema: schema}
}

func (m *Postgres) Migrate(tx *sql.Tx, metaModels []MetaModel) {
	for _, model := range metaModels {
		m.migrate(tx, &model)
	}
}

func (m *Postgres) DropTable(tx *sql.Tx, name string) {
	name = naming.Get(name)

	if _, err := db.Exec(`DROP TABLE "` + name + `"`); err != nil {
		panic(err)
	}
}

func (m *Postgres) migrate(tx *sql.Tx, model *MetaModel) {
	if !m.tableExists(model.ModelName) {
		m.createTable(tx, model)
	} else {
		m.updateTable(tx, model)
	}
}

func (m *Postgres) tableExists(name string) bool {
	name = naming.Get(name)

	rows, err := testdb.Query(`SELECT EXISTS (SELECT 1
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

func (m *Postgres) createTable(tx *sql.Tx, model *MetaModel) {
	name := naming.Get(model.ModelName)
	sql := `CREATE TABLE "` + name + `" (` + m.getColumns(model) + `)`

	// create sequences if required
	for _, seq := range m.sequences {
		if _, err := tx.Exec(seq); err != nil {
			panic(err)
		}
	}

	// create table
	if _, err := tx.Exec(sql); err != nil {
		panic(err)
	}

	// alter sequence if required
	for _, seq := range m.alterSeq {
		if _, err := tx.Exec(seq); err != nil {
			panic(err)
		}
	}

	// reset
	m.sequences = make([]string, 0)
	m.alterSeq = make([]string, 0)
}

func (m *Postgres) updateTable(tx *sql.Tx, model *MetaModel) {

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
		} else if key == "seq" {
			m.addSequence(modelName, field.Name, value)
		} else if value == "id" {
			// id is a shortcut for seq + default
			m.addSequence(modelName, field.Name, "1,1,-,-,1")
			tags[1] = "DEFAULT nextval('" + m.getSequenceName(modelName, field.Name) + "'::regclass)"
		} else if value == "pk" {
			tags[3] = "PRIMARY KEY"
		} else if value == "unique" {
			tags[4] = "UNIQUE"
		}
	}

	return strings.Join(tags, " ")
}

func (m *Postgres) addSequence(modelName, columnName, info string) {
	// create sequence
	infos := strings.Split(info, ",")

	if len(infos) != 5 {
		panic("5 arguments must be specified for seq: start, increment, min, max, cache")
	}

	name := m.getSequenceName(modelName, columnName)
	seq := `CREATE SEQUENCE "` + name + `"
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

	m.sequences = append(m.sequences, seq)

	// create owned by table
	modelName = naming.Get(modelName)
	columnName = naming.Get(columnName)
	alterSeq := `ALTER SEQUENCE "` + name + `" OWNED BY "` + modelName + `"."` + columnName + `"`
	m.alterSeq = append(m.alterSeq, alterSeq)
}

func (m *Postgres) getSequenceName(modelName, columnName string) string {
	modelName = naming.Get(modelName)
	columnName = naming.Get(columnName)
	return modelName + "_" + columnName + "_seq"
}
