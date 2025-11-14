package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/mohsen104/web-api/api/handlers"
)

func User(r *gin.RouterGroup) {
	handler := handlers.NewUserHandler()

	r.GET("/", handler.User)
	r.GET("/:id", handler.User)
}
