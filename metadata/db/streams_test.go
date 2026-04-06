package db_test

import (
	"testing"

	"github.com/astria-tv/astria-server/metadata/app"
	"github.com/astria-tv/astria-server/metadata/db"
)

func TestBeforeCreate(t *testing.T) {
	app.NewMDContext(db.DatabaseOptions{
		Connection: db.InMemory,
	}, nil)
	stream := db.Stream{Codecs: "test"}
	db.CreateStream(&stream)
	if stream.UUID == "" {
		t.Errorf("Stream was created without a UUID\n")
	}
}
