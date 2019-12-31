package user

import "github.com/gin-gonic/gin"

func RegisterHandler(r *gin.Engine, service Service) {
	res := resource{service}

	userRoute := r.Group("user")
	{
		userRoute.POST("/register", res.create)
	}
}

type resource struct {
	service Service
}

func (r resource) create(c *gin.Context) {
	var dto CreateUserDto
	c.ShouldBind(&dto)
	r.service.Create(c, dto)
	return
}
