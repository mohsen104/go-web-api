package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/mohsen104/web-api/api/handlers"
)

func TestRouter(r *gin.RouterGroup) {
	h := handlers.NewTestHandler()

	r.GET("/", h.Test)
	r.GET("/:id", h.TestId)
	r.GET("/header", h.Header)
	r.GET("/query", h.Query)
	r.GET("/create", h.Create)
}
