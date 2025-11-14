package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/mohsen104/web-api/api/helper"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// User godoc
// @Summary Get user
// @Description Get user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} helper.BaseHttpResponse "Success"
// @Failure 400 {object} helper.BaseHttpResponse "Failed"
// @Router /v1/users [get]
func (u *UserHandler) User(c *gin.Context) {
	c.JSON(200,
		helper.GenerateBaseResponse("hello world", true, 0),
	)
}

// User godoc
// @Summary Get user by id
// @Description Get user by id
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} helper.BaseHttpResponse "Success"
// @Failure 400 {object} helper.BaseHttpResponse "Failed"
// @Router /v1/users/{id} [get]
func (u *UserHandler) UserById(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200,
		helper.GenerateBaseResponse("hello world "+id, true, 0),
	)
}
