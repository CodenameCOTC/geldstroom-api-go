package main

import (
	"flag"

	"github.com/novaladip/geldstroom-api-go/server"
)

func main() {
	addr := flag.String("addr", ":4000", "Http network address")
	dsn := flag.String("dsn", "root@/geldstroom?parseTime=true", "MySQL data source name")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	flag.Parse()

	server := &server.Server{
		Dsn:    *dsn,
		Addr:   *addr,
		Secret: *secret,
	}

	server.Start()

}
