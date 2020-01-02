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
	"github.com/novaladip/geldstroom-api-go/core/transaction"
	"github.com/novaladip/geldstroom-api-go/core/user"
	"github.com/novaladip/geldstroom-api-go/pkg/database"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error load .env file")
	}

	config.LoadKey()

	db, err := database.OpenDB(config.ConfigKey.DB_DSN)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer db.Close()

	srv := &http.Server{
		Addr:         config.ConfigKey.ADDR,
		Handler:      buildHanlder(db, &config.ConfigKey),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Server listening on PORT: " + config.ConfigKey.ADDR)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func buildHanlder(db *sql.DB, key *config.Key) http.Handler {
	router := gin.Default()

	user.RegisterHandler(router,
		user.NewService(user.NewRepository(db)),
	)
	transaction.RegisterHandler(router, db)

	return router
}
