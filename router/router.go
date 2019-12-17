package router

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/novaladip/geldstroom-api-go/auth"
)

type Router struct {
	DB *sql.DB
	R  *gin.Engine
}

// Initializing routes
func (r Router) Init() {
	auth := &auth.AuthDb{
		Db: r.DB,
	}

	authRoutes := r.R.Group("/auth")
	{
		authRoutes.POST("/login", auth.Login)
		authRoutes.POST("/register", auth.Register)
	}
}
