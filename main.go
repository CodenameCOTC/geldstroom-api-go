package main

import (
	"github.com/joho/godotenv"
	"github.com/novaladip/geldstroom-api-go/logger"
	"github.com/novaladip/geldstroom-api-go/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.ErrorLog.Fatal("Error loading .env file")
	}

	server, db := server.New()
	defer db.Close()

	logger.InfoLog.Println("Server starting to listen & serve")

	err = server.ListenAndServe()

	if err != nil {
		logger.ErrorLog.Fatal(err)
	}

}
