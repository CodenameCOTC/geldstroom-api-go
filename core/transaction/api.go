package transaction

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/novaladip/geldstroom-api-go/core/auth"
	"github.com/novaladip/geldstroom-api-go/core/entity"
	"github.com/novaladip/geldstroom-api-go/pkg/getrange"

	errorsresponse "github.com/novaladip/geldstroom-api-go/core/errors"
)

func RegisterHandler(r *gin.Engine, db *sql.DB, service Service) {
	res := resource{service}

	authMiddleare := auth.NewMiddleware(auth.NewRepository(db))

	transactionRoutes := r.Group("/transaction")
	transactionRoutes.Use(authMiddleare.AuthGuard())
	{
		transactionRoutes.GET("/", res.get)
		transactionRoutes.POST("/", res.create)
		transactionRoutes.GET("/:id", res.findOneById)
		transactionRoutes.DELETE("/:id", res.deleteOneById)
		transactionRoutes.PUT("/:id", res.updateOneById)
	}
}

type resource struct {
	service Service
}

func (r resource) get(c *gin.Context) {
	user, _ := c.MustGet("JwtPayload").(entity.JwtPayload)
	dateRange := strings.ToUpper(c.Query("range"))
	date := c.Query("date")

	dr, err := getrange.GetRange(date, dateRange)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorsresponse.InvalidQuery(ErrInvalidQueryCode, err))
		return
	}

	t, err := r.service.Get(*dr, user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorsresponse.InternalServerError(""))
		return
	}

	c.JSON(http.StatusOK, t)
}

func (r resource) create(c *gin.Context) {
	var dto CreateDto
	_ = c.ShouldBind(&dto)
	user, _ := c.MustGet("JwtPayload").(entity.JwtPayload)

	validate := dto.Validate()

	if !validate.IsValid {
		c.JSON(http.StatusBadRequest, errorsresponse.
			ValidationError(ErrValidationFailedCode,
				ErrValidationFailed,
				validate.Error),
		)
		return
	}

	t, err := r.service.Create(entity.Transaction{
		Id:          entity.GenerateID(),
		Amount:      dto.Amount,
		Description: dto.Description,
		Category:    strings.ToUpper(dto.Category),
		Type:        strings.ToUpper(dto.Type),
		UserId:      user.Id,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorsresponse.InternalServerError(""))
	}

	c.JSON(http.StatusOK, t)
}

func (r resource) findOneById(c *gin.Context) {
	user, _ := c.MustGet("JwtPayload").(entity.JwtPayload)
	tId := c.Param("id")

	t, err := r.service.FindOneById(tId, user.Id)
	if err != nil {
		if errors.Is(err, ErrTransactionNotFound) {
			c.JSON(http.StatusNotFound, errorsresponse.NotFound(fmt.Sprintf("Transaction with ID: %v is not found", tId)))
			return
		}
		c.JSON(http.StatusInternalServerError, errorsresponse.InternalServerError(""))
		return
	}

	c.JSON(http.StatusOK, t)
}

func (r resource) deleteOneById(c *gin.Context) {
	user, _ := c.MustGet("JwtPayload").(entity.JwtPayload)
	tId := c.Param("id")

	if err := r.service.DeleteOneById(tId, user.Id); err != nil {
		if errors.Is(err, ErrTransactionNotFound) {
			c.JSON(http.StatusNotFound, errorsresponse.NotFound(fmt.Sprintf("Transaction with ID: %v is not found", tId)))
			return
		}
		c.JSON(http.StatusInternalServerError, errorsresponse.InternalServerError(""))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Transaction with ID: %v has ben deleted", tId),
	})
}

func (r resource) updateOneById(c *gin.Context) {
	var dto UpdateDto
	user, _ := c.MustGet("JwtPayload").(entity.JwtPayload)
	tId := c.Param("id")

	_ = c.ShouldBind(&dto)

	validate := dto.Validate()
	if !validate.IsValid {
		c.JSON(http.StatusBadRequest, errorsresponse.
			ValidationError(ErrValidationFailedCode,
				ErrValidationFailed,
				validate.Error),
		)
		return
	}

	t, err := r.service.UpdateOneById(tId, user.Id, dto)
	fmt.Println(err)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorsresponse.NotFound(fmt.Sprintf("Transaction with ID: %v is not found", tId)))
			return
		}
		c.JSON(http.StatusInternalServerError, errorsresponse.InternalServerError(""))
		return
	}

	c.JSON(http.StatusOK, t)

}
