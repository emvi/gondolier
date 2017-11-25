package gondolier

import (
	"database/sql"
	"log"
	"strings"
)

// Migrator for Postgres databases.
// You can use the following options to configure your data model:
//
//  // The type must be the database type.
//  type:database type
//  // Sets the column as primary key.
//  pk/primary key
//  // Creates and sets a sequence with given parameters for the column.
//  seq:start,increment,minvalue,maxvalue,cache
//  // Sets the default value for column, strings must be escaped.
//  // next(seq) refers to the sequences assign for this column (using seq:...).
//  default:default value/next(seq)
//  // Sets not null constraint for column.
//  not null/notnull
//  // Optional. Drops not null constraint if set for column. Not null is also dropped if not null is not set.
//  null
//  // Sets unique constraint for column.
//  unique
//  // Shortcut for primary key, not null, seq:1,1,-,-,1 and default:next(seq).
//  id
//  // Sets foreign key constraint for column.
//  // It refers to the given model and column.
//  // Example: fk:MyModel.Id
//  fk/foreign key:Model.Column
type Postgres struct {
	Schema      string
	DropColumns bool
	Log         bool

	createSeq []string
	alterSeq  []string
	createFK  []string
	dropFK    []string
	alterPK   string
}

func (m *Postgres) Migrate(metaModels []MetaModel) {
	// create or update table
	for _, model := range metaModels {
		m.migrate(&model)
	}

	// create foreign keys
	for _, fk := range m.createFK {
		m.exec(fk)
	}

	// drop foreign keys
	for _, fk := range m.dropFK {
		m.exec(fk)
	}

	// reset
	m.createFK = make([]string, 0)
	m.dropFK = make([]string, 0)
}

func (m *Postgres) DropTable(name string) {
	name = naming.Get(name)
	m.exec(`DROP TABLE IF EXISTS "` + name + `"`)
}

func (m *Postgres) migrate(model *MetaModel) {
	if !m.tableExists(model.ModelName) {
		m.createTable(model)
	} else {
		m.updateTable(model)

		if m.DropColumns {
			m.dropColumns(model)
		}
	}
}

func (m *Postgres) tableExists(name string) bool {
	name = naming.Get(name)

	rows, err := db.Query(`SELECT EXISTS (SELECT 1
	   FROM information_schema.tables
	   WHERE table_schema = $1
	   AND table_name = $2)`, m.Schema, name)

	return m.scanBool(rows, err)
}

func (m *Postgres) columnExists(tableName, columnName string) bool {
	tableName = naming.Get(tableName)
	columnName = naming.Get(columnName)

	rows, err := db.Query(`SELECT EXISTS (SELECT 1
	   FROM information_schema.columns
	   WHERE table_schema = $1
	   AND table_name = $2
	   AND column_name = $3)`, m.Schema, tableName, columnName)

	return m.scanBool(rows, err)
}

func (m *Postgres) sequenceExists(name string) bool {
	name = naming.Get(name)

	rows, err := db.Query(`SELECT EXISTS (SELECT 1
	   FROM pg_class
	   WHERE relkind = 'S'
	   AND oid::regclass::text = quote_ident($1))`, name)

	return m.scanBool(rows, err)
}

func (m *Postgres) foreignKeyExists(tableName, fkName string) bool {
	tableName = naming.Get(tableName)
	fkName = naming.Get(fkName)

	rows, err := db.Query(`SELECT EXISTS (SELECT 1
		FROM information_schema.table_constraints
		WHERE table_schema = $1
		AND constraint_name = $2
		AND table_name = $3)`, m.Schema, fkName, tableName)

	return m.scanBool(rows, err)
}

func (m *Postgres) isNullable(tableName, columnName string) bool {
	tableName = naming.Get(tableName)
	columnName = naming.Get(columnName)

	rows, err := db.Query(`SELECT is_nullable::boolean
		FROM information_schema.columns
		WHERE table_schema = $1
		AND column_name = $2
		AND table_name = $3`, m.Schema, columnName, tableName)

	return m.scanBool(rows, err)
}

func (m *Postgres) constraintExists(name string) bool {
	name = naming.Get(name)

	rows, err := db.Query(`SELECT EXISTS (SELECT 1
		FROM pg_constraint WHERE conname = $1)`, name)

	return m.scanBool(rows, err)
}

func (m *Postgres) scanBool(rows *sql.Rows, err error) bool {
	if err != nil {
		panic(err)
	}

	var exists bool
	rows.Next()

	if err := rows.Scan(&exists); err != nil {
		panic(err)
	}

	m.closeRows(rows)
	return exists
}

func (m *Postgres) getColumnNames(tableName string) []string {
	tableName = naming.Get(tableName)

	rows, err := db.Query(`SELECT column_name
		FROM information_schema.columns
		WHERE table_schema = $1
		AND table_name = $2`, m.Schema, tableName)

	if err != nil {
		panic(err)
	}

	names := make([]string, 0)

	for rows.Next() {
		var name string

		if err := rows.Scan(&name); err != nil {
			panic(err)
		}

		names = append(names, name)
	}

	m.closeRows(rows)
	return names
}

func (m *Postgres) getColumnType(tableName, columnName string) string {
	tableName = naming.Get(tableName)
	columnName = naming.Get(columnName)

	rows, err := db.Query(`SELECT data_type FROM information_schema.columns
		WHERE table_name = $1 AND column_name = $2`, tableName, columnName)

	if err != nil {
		panic(err)
	}

	var typeName string
	rows.Next()

	if err := rows.Scan(&typeName); err != nil {
		panic(err)
	}

	m.closeRows(rows)
	return typeName
}

func (m *Postgres) getConstraintName(name string) string {
	name = naming.Get(name)

	rows, err := db.Query(`SELECT conname
		FROM pg_constraint WHERE conname LIKE $1`, name)

	if err != nil {
		panic(err)
	}

	var constraintName string
	one := false

	for rows.Next() {
		if err := rows.Scan(&constraintName); err != nil {
			panic(err)
		}

		if one {
			panic("No distinct constraint found for name '" + name + "'")
		}

		one = true
	}

	m.closeRows(rows)
	return constraintName
}

func (m *Postgres) closeRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		panic(err)
	}
}

func (m *Postgres) createTable(model *MetaModel) {
	name := naming.Get(model.ModelName)
	sql := `CREATE TABLE IF NOT EXISTS "` + name + `" (` + m.getColumns(model) + `)`

	// create sequences if required
	for _, seq := range m.createSeq {
		m.exec(seq)
	}

	// create table
	m.exec(sql)

	// alter sequence if required
	for _, seq := range m.alterSeq {
		m.exec(seq)
	}

	// alter primary key if required
	if m.alterPK != "" {
		m.exec(m.alterPK)
	}

	// reset
	m.createSeq = make([]string, 0)
	m.alterSeq = make([]string, 0)
	m.alterPK = ""
}

func (m *Postgres) updateTable(model *MetaModel) {
	for _, field := range model.Fields {
		if m.columnExists(model.ModelName, field.Name) {
			// update existing column
			m.updateColumn(model, &field)
		} else {
			// create new column
			tableName := naming.Get(model.ModelName)
			columnName := naming.Get(field.Name)
			query := `ALTER TABLE "` + tableName + `" ADD COLUMN "` + columnName + `" ` + m.getTags(tableName, &field)
			m.exec(query)
		}
	}
}

func (m *Postgres) updateColumn(model *MetaModel, field *MetaField) {
	tableName := naming.Get(model.ModelName)
	columnName := naming.Get(field.Name)
	notnull, isId, pk, unique := false, false, false, false
	defaultValue, seq, fk := "", "", ""

	for _, tag := range field.Tags {
		key := strings.ToLower(tag.Name)
		value := strings.ToLower(tag.Value)

		if key == "type" {
			m.updateColumnType(tableName, columnName, value)
		} else if value == "notnull" || value == "not null" {
			notnull = true
		} else if value == "null" {
			notnull = false
		} else if key == "default" {
			defaultValue = value
		} else if value == "id" {
			notnull = true
			isId = true
			pk = true
		} else if value == "pk" || value == "primary key" {
			pk = true
		} else if value == "unique" {
			unique = true
		} else if key == "seq" || key == "sequence" {
			seq = value
		} else if key == "fk" || key == "foreign key" {
			fk = tag.Value
		}
	}

	m.updateColumnSeq(tableName, columnName, seq, isId)
	m.updateColumnPK(tableName, columnName, pk)
	m.updateColumnUnique(tableName, columnName, unique)
	m.updateColumnNotNull(tableName, columnName, notnull)
	m.updateColumnDefault(tableName, columnName, defaultValue, isId)
	m.updateColumnFk(tableName, columnName, fk)
}

func (m *Postgres) updateColumnType(tableName, columnName, newtype string) {
	istype := m.getColumnType(tableName, columnName)

	if istype != newtype {
		query := `ALTER TABLE "` + tableName + `" ALTER COLUMN "` + columnName + `"
					TYPE ` + newtype
		m.exec(query)
	}
}

func (m *Postgres) updateColumnNotNull(tableName, columnName string, notnull bool) {
	query := `ALTER TABLE "` + tableName + `" ALTER COLUMN "` + columnName + `"`

	if notnull {
		query += " SET NOT NULL"
	} else {
		query += " DROP NOT NULL"
	}

	m.exec(query)
}

func (m *Postgres) updateColumnDefault(tableName, columnName, value string, isId bool) {
	query := ""

	if value != "" || isId {
		// set default
		if isId {
			m.addSequence(tableName, columnName, "1,1,-,-,1")
			m.exec(m.createSeq[0])
			m.exec(m.alterSeq[0])
			m.createSeq = make([]string, 0)
			m.alterSeq = make([]string, 0)
			query = `ALTER TABLE "` + tableName + `" ALTER COLUMN "` + columnName + `" SET DEFAULT nextval('` + m.getSequenceName(tableName, columnName) + `'::regclass)`
		} else {
			if value == "nextval(seq)" {
				value = "nextval('" + m.getSequenceName(tableName, columnName) + "'::regclass)"
			}

			query = `ALTER TABLE "` + tableName + `" ALTER COLUMN "` + columnName + `" SET DEFAULT ` + value
		}
	} else {
		// drop default
		query = `ALTER TABLE "` + tableName + `" ALTER COLUMN "` + columnName + `" DROP DEFAULT`
	}

	m.exec(query)
}

func (m *Postgres) updateColumnPK(tableName, columnName string, pk bool) {
	pkName := m.getPrimaryKeyName(tableName, columnName)

	if !pk && m.constraintExists(pkName) {
		m.exec(`ALTER TABLE "` + tableName + `" DROP CONSTRAINT IF EXISTS "` + pkName + `"`)
	} else if pk && !m.constraintExists(pkName) {
		m.exec(`ALTER TABLE "` + tableName + `" ADD PRIMARY KEY ("` + columnName + `")`)
	}
}

func (m *Postgres) updateColumnUnique(tableName, columnName string, unique bool) {
	query := ""
	constraintName := m.getUniqueName(tableName, columnName)

	if unique && !m.constraintExists(constraintName) {
		query = `ALTER TABLE "` + tableName + `" ADD CONSTRAINT "` + constraintName + `" UNIQUE ("` + columnName + `")`
	} else if !unique && m.constraintExists(constraintName) {
		query = `ALTER TABLE "` + tableName + `" DROP CONSTRAINT IF EXISTS "` + constraintName + `"`
	}

	m.exec(query)
}

func (m *Postgres) updateColumnSeq(tableName, columnName, seq string, isId bool) {
	if isId {
		return
	}

	seqName := m.getSequenceName(tableName, columnName)

	if seq != "" && !m.sequenceExists(seqName) {
		// create sequence
		m.addSequence(tableName, columnName, seq)
		m.exec(m.createSeq[0])
		m.exec(m.alterSeq[0])
		m.createSeq = make([]string, 0)
		m.alterSeq = make([]string, 0)
	} else if seq == "" && m.sequenceExists(seqName) {
		// drop sequence
		query := `DROP SEQUENCE IF EXISTS "` + seqName + `" CASCADE`
		m.exec(query)
	}
}

func (m *Postgres) updateColumnFk(tableName, columnName, fk string) {
	// read existing fk
	refTableName, refColumnName := m.getForeignKeyInfo(tableName, fk)
	fkName := m.getForeignKeyName(tableName, columnName, refTableName, refColumnName)
	existingFk := m.getConstraintName(tableName + "_" + columnName + "_%_fk")

	if fkName != existingFk {
		// drop on change or when it was removed if exists
		if existingFk != "" {
			m.dropFK = append(m.dropFK, `ALTER TABLE "`+tableName+`" DROP CONSTRAINT IF EXISTS "`+existingFk+`"`)
		}

		// create new
		if fk != "" {
			m.addForeignKey(tableName, columnName, fk)
		}
	}
}

// Drops all columns that are no longer needed.
func (m *Postgres) dropColumns(model *MetaModel) {
	tableName := naming.Get(model.ModelName)
	columns := m.getColumnNames(model.ModelName)

	for _, column := range columns {
		if !m.fieldsContainsColumn(model.Fields, column) {
			query := `ALTER TABLE "` + tableName + `"
				DROP COLUMN IF EXISTS "` + column + `"`
			m.exec(query)
		}
	}
}

func (m *Postgres) fieldsContainsColumn(fields []MetaField, column string) bool {
	for _, field := range fields {
		if naming.Get(field.Name) == column {
			return true
		}
	}

	return false
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
		} else if key == "seq" || key == "sequence" {
			m.addSequence(modelName, field.Name, value)
		} else if value == "id" {
			// id is a shortcut for seq + default + pk
			m.addSequence(modelName, field.Name, "1,1,-,-,1")
			tags[1] = "DEFAULT nextval('" + m.getSequenceName(modelName, field.Name) + "'::regclass)"
			tags[3] = "PRIMARY KEY"
			m.alterPrimaryKey(modelName, field.Name)
		} else if value == "pk" || value == "primary key" {
			tags[3] = "PRIMARY KEY"
			m.alterPrimaryKey(modelName, field.Name)
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

func (m *Postgres) alterPrimaryKey(modelName, columnName string) {
	tableName := naming.Get(modelName)
	m.alterPK = `ALTER TABLE "` + tableName + `"
		RENAME CONSTRAINT "` + tableName + `_pkey"
		TO "` + m.getPrimaryKeyName(modelName, columnName) + `"`
}

func (m *Postgres) getSequenceName(modelName, columnName string) string {
	modelName = naming.Get(modelName)
	columnName = naming.Get(columnName)
	return modelName + "_" + columnName + "_seq"
}

func (m *Postgres) addForeignKey(modelName, columnName, info string) {
	refTableName, refColumnName := m.getForeignKeyInfo(modelName, info)
	tableName := naming.Get(modelName)
	columnName = naming.Get(columnName)
	fkName := m.getForeignKeyName(modelName, columnName, refTableName, refColumnName)
	alterFk := `ALTER TABLE "` + tableName + `"
		ADD CONSTRAINT "` + fkName + `"
		FOREIGN KEY ("` + columnName + `")
		REFERENCES "` + refTableName + `"("` + refColumnName + `")`
	m.createFK = append(m.createFK, alterFk)
}

func (m *Postgres) getForeignKeyInfo(modelName, info string) (string, string) {
	if info == "" {
		return "", ""
	}

	infos := strings.Split(info, ".")

	if len(infos) != 2 {
		panic("Two arguments must be specified for fk in model '" + modelName + "': ReferencedModel.ReferencedAttribute")
	}

	return naming.Get(infos[0]), naming.Get(infos[1])
}

func (m *Postgres) getForeignKeyName(modelName, columnName, refObjName, refColumnName string) string {
	modelName = naming.Get(modelName)
	refObjName = naming.Get(refObjName)
	refColumnName = naming.Get(refColumnName)
	return modelName + "_" + columnName + "_" + refObjName + "_" + refColumnName + "_fk"
}

func (m *Postgres) getPrimaryKeyName(modelName, columnName string) string {
	modelName = naming.Get(modelName)
	columnName = naming.Get(columnName)
	return modelName + "_" + columnName + "_pkey"
}

func (m *Postgres) getUniqueName(modelName, columnName string) string {
	modelName = naming.Get(modelName)
	columnName = naming.Get(columnName)
	return modelName + "_" + columnName + "_key"
}

func (m *Postgres) exec(query string) {
	if m.Log {
		log.Println(query)
	}

	if _, err := db.Exec(query); err != nil {
		panic(err)
	}
}
