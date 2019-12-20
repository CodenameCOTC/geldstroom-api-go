package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/novaladip/geldstroom-api-go/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	server, db := server.New()
	defer db.Close()

	err = server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}

}
