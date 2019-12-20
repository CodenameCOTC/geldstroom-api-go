package router

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/novaladip/geldstroom-api-go/auth"
	"github.com/novaladip/geldstroom-api-go/middleware"
	"github.com/novaladip/geldstroom-api-go/transaction"
)

type Router struct {
	DB *sql.DB
	R  *gin.Engine
}

// Initializing routes
func (r Router) Init() {

	guard := &middleware.Guard{
		Db: r.DB,
	}

	authHandler := &auth.Handler{
		Db: r.DB,
	}
	authRoutes := r.R.Group("/auth")
	{
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/register", authHandler.Register)
	}

	transactionHandler := &transaction.Handler{
		Db: r.DB,
	}

	transactionRoutes := r.R.Group("/transaction")
	transactionRoutes.Use(guard.AuthGuard())
	{
		transactionRoutes.GET("/", transactionHandler.GetTransactions)
		transactionRoutes.POST("/", transactionHandler.Create)
		transactionRoutes.PUT("/:id", transactionHandler.Update)
		transactionRoutes.DELETE("/:id", transactionHandler.Delete)
	}

}
