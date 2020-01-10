package transaction

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/novaladip/geldstroom-api-go/core/auth"
	"github.com/novaladip/geldstroom-api-go/pkg/entity"
	"github.com/novaladip/geldstroom-api-go/pkg/getrange"
	"github.com/novaladip/geldstroom-api-go/pkg/pagination"

	errorsresponse "github.com/novaladip/geldstroom-api-go/pkg/errors"
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
	user := entity.JwtPayloadFromRequest(c)
	p := pagination.NewFromRequest(c)
	dr, err := getrange.NewFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorsresponse.InvalidQuery(ErrInvalidQueryCode, err))
		return
	}

	t, count, err := r.service.Get(*dr, p.Page, p.PerPage, user.Id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorsresponse.InternalServerError(""))
		return
	}

	p.Items = t
	p.TotalCount = count
	p.PageCount = int(math.Ceil(float64(count) / float64(p.PerPage)))

	c.JSON(http.StatusOK, p)
}

func (r resource) create(c *gin.Context) {
	dto := NewCreateDtoFromRequest(c)
	user := entity.JwtPayloadFromRequest(c)
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
	tId := c.Param("id")
	user := entity.JwtPayloadFromRequest(c)

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
	tId := c.Param("id")
	user := entity.JwtPayloadFromRequest(c)

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
	tId := c.Param("id")
	user := entity.JwtPayloadFromRequest(c)
	dto := NewUpdateDtoFromRequest(c)

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
