package gondolier

import (
	"database/sql"
)

type Postgres struct{}

func (t *Postgres) Migrate(tx *sql.Tx, metaModels []MetaModel) {

}

func (t *Postgres) DropTable(tx *sql.Tx, name string) {
	name = naming.Get(name)

	if _, err := db.Exec(`DROP TABLE "` + name + `"`); err != nil {
		panic(err)
	}
}
