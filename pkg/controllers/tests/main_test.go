package controllers

import (
	"os"
	"simplebank/pkg/connection"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func TestMain(m *testing.M) {
	db := connection.OpenConnection()

	DB = db

	os.Exit(m.Run())
}
