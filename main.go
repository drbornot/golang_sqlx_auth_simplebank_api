package main

import (
	"log"
	"simplebank/pkg/connection"
)

func main() {
	// ctx := context.Background()

	db := connection.OpenConnection()
	defer db.Close()

	log.Println("Connection Closed!")
}
