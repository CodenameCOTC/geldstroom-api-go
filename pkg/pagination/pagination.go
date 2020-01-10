package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pages struct {
	Page       int         `json:"page"`
	PerPage    int         `json:"perPage"`
	TotalCount int         `json:"totalCount"`
	PageCount  int         `json:"pageCount"`
	Items      interface{} `json:"items"`
}

var (
	DefaultPage    = 1
	DefaultPerPage = 10
	PageVar        = "page"
	PerPageVar     = "perPage"
)

/**
.... HttpRequest -> NewFromRequst -> Repository
*/
func NewFromRequest(c *gin.Context) *Pages {
	page := parseInt(c.Query(PageVar), DefaultPage, 0)
	perPage := parseInt(c.Query(PerPageVar), DefaultPerPage, 0)

	p := &Pages{
		Page:    page,
		PerPage: perPage,
	}

	return p
}

func parseInt(value string, defaultValue, limit int) int {
	if value == "" {
		return defaultValue
	}
	if result, err := strconv.Atoi(value); err == nil {
		if result < limit {
			return defaultValue
		}
		return result
	}
	return defaultValue
}
