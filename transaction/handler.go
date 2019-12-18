package transaction

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/novaladip/geldstroom-api-go/auth"
	"github.com/novaladip/geldstroom-api-go/helper"
)

type Handler struct {
	Db *sql.DB
}

func (h *Handler) Create(c *gin.Context) {
	var dto insertDto
	user, _ := c.MustGet("JwtPayload").(auth.JwtPayload)

	c.ShouldBind(&dto)

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

func (h *Handler) Update(c *gin.Context) {
	user, _ := c.MustGet("JwtPayload").(auth.JwtPayload)
	tId := c.Param("id")
	var dto updateDto
	c.ShouldBind(&dto)
	vud := newValidateUpdateDto(&dto)

	if ok := vud.validate(); !ok {
		c.JSON(http.StatusBadRequest, &helper.ErrorResponse{
			Message:   "Field Error",
			ErrorCode: "TRANSACTION_0002",
			Error:     vud.error,
		})
		return
	}

	t, err := h.update(tId, dto, user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, t)
}

func (h *Handler) Delete(c *gin.Context) {
	user, _ := c.MustGet("JwtPayload").(auth.JwtPayload)
	tId := c.Param("id")

	if err := h.delete(tId, user.Id); err != nil {
		if errors.Is(err, helper.ErrSqlNoRow) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Transaction not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, &helper.InternalServerError)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
