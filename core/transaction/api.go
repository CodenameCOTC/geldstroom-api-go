package transaction

import "github.com/gin-gonic/gin"

import "net/http"

import "github.com/novaladip/geldstroom-api-go/core/auth"

import "database/sql"

import "github.com/novaladip/geldstroom-api-go/core/entity"

func RegisterHandler(r *gin.Engine, db *sql.DB) {
	authMiddleare := auth.NewMiddleware(auth.NewRepository(db))

	transactionRoutes := r.Group("/transaction")
	transactionRoutes.Use(authMiddleare.AuthGuard())
	{
		transactionRoutes.GET("/", get)
	}
}

func get(c *gin.Context) {
	user, _ := c.MustGet("JwtPayload").(entity.JwtPayload)
	c.JSON(http.StatusOK, user)
}
