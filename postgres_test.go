package gondolier

import (
	"testing"
)

type testUser struct {
	Id      uint64 `gondolier:"type:bigint;pk;notnull"`
	Name    string `gondolier:"type:character varying(255);notnull"`
	Age     uint   `gondolier:"type:integer";notnull`
	Picture uint64 `gondolier:"type:bigint;fk:testPicture;null"`
}

type testPicture struct {
	Id       uint64 `gondolier:"type:bigint;pk;notnull"`
	FileName string `gondolier:"type:character varying(255);notnull"`
	User     uint64 `gondolier:"type:bigint;fk:testUser;null"`
}

type testPost struct {
	Id      uint64 `gondolier:"type:bigint;pk;notnull"`
	Post    string `gondolier:"type:character varying(255);notnull"`
	User    uint64 `gondolier:"type:bigint;fk:testUser;notnull"`
	Picture uint64 `gondolier:"type:bigint;fk:testPicture;null"`
}

func TestPostgresMigrator(t *testing.T) {
	Use(&Postgres{})
	Model(testUser{}, testPicture{}, testPost{})
	Migrate()
}
