package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/novaladip/geldstroom-api-go/config"
	"github.com/novaladip/geldstroom-api-go/database"
	"github.com/novaladip/geldstroom-api-go/router"
)

func New() (*http.Server, *sql.DB) {
	fmt.Println("Initializing server...")

	key := config.GetKey()

	db, err := database.OpenDB(key.DB_DSN)
	if err != nil {
		log.Fatal(err.Error())
	}

	r := gin.Default()

	route := &router.Router{DB: db, R: r}

	route.Init()

	srv := &http.Server{
		Addr:         key.ADDR,
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return srv, db
}
