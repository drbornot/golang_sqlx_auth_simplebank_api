package main

import (
	"log"
	"simplebank/api"
	"simplebank/pkg/connection"
)

var addressServer = "0.0.0.0:8080"

func main() {

	db := connection.OpenConnection()
	defer db.Close()

	server := api.NewServer(db)

	err := server.Start(addressServer)
	if err != nil {
		log.Fatal("cannot start server", err)
	}

	log.Println("Connection Closed!")
}
