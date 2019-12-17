package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/novaladip/geldstroom-api-go/database"
	"github.com/novaladip/geldstroom-api-go/router"
)

type Server struct {
	Dsn    string
	Addr   string
	Secret string
}

func (s Server) Start() {
	fmt.Println("Initializing server...")

	db, err := database.OpenDB(s.Dsn)
	if err != nil {
		log.Fatal(err.Error())
	}

	r := gin.Default()

	defer db.Close()

	route := &router.Router{DB: db, R: r, Secret: s.Secret}

	route.Init()

	srv := &http.Server{
		Addr:         s.Addr,
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
