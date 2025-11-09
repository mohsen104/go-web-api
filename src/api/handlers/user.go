package handlers

import "github.com/gin-gonic/gin"

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (u *UserHandler) User(c *gin.Context) {
	c.JSON(200, "hello world")
	return
}
