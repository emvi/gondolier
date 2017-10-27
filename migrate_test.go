package model

import (
	"testing"
)

func TestUse(t *testing.T) {
	Use("Postgres")

	if database != Postgres {
		t.Fatal("Postgres must be selected")
	}
}

func TestUseUnknownDb(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Fatal("Calling Use with unknown database must panic")
		}
	}()

	Use("unknown")
}
