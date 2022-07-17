package connection

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func OpenConnection() *sqlx.DB {
	db, err := sqlx.Open("postgres", "user=root password=secret dbname=simple_bank sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Connected!")
	return db
}
