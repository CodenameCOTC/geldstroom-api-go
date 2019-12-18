package transaction

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/novaladip/geldstroom-api-go/auth"
	"github.com/novaladip/geldstroom-api-go/helper"
)

type Handler struct {
	Db *sql.DB
}

func (h *Handler) Create(c *gin.Context) {
	var dto InsertDto
	user, _ := c.MustGet("JwtPayload").(auth.JwtPayload)

	if err := c.ShouldBind(&dto); err != nil {
		c.JSON(http.StatusBadRequest, helper.Unauthorized)
		return
	}

	vid := newValidateInsertDto(&dto)

	if ok := vid.validate(); !ok {
		c.JSON(http.StatusBadRequest, &helper.ErrorResponse{
			Message:   errInsert.Error(),
			ErrorCode: errInsertCode,
			Error:     vid.error,
		})
		return
	}

	t, err := h.insert(dto, user.Id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.InternalServerError)
		return
	}

	c.JSON(http.StatusCreated, t)

}

func (h *Handler) GetTransactions(c *gin.Context) {
	user, ok := c.MustGet("JwtPayload").(auth.JwtPayload)

	if !ok {
		c.JSON(http.StatusUnauthorized, user)
		return
	}

	t, err := h.getTransaction(&user.Id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, &t)

}

func (h *Handler) Update(c *gin.Context) {}

func (h *Handler) Delete(c *gin.Context) {}
