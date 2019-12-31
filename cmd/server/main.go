package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/novaladip/geldstroom-api-go/core/config"
	"github.com/novaladip/geldstroom-api-go/pkg/database"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error load .env file")
	}

	key := config.GetKey()

	db, err := database.OpenDB(key.DB_DSN)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer db.Close()

	srv := &http.Server{
		Addr:         key.ADDR,
		Handler:      buildHanlder(db, key),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Server listening on PORT: " + key.ADDR)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func buildHanlder(db *sql.DB, key *config.Key) http.Handler {
	router := gin.Default()

	return router
}
