package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mohsen104/web-api/api/routers"
	"github.com/mohsen104/web-api/config"
)

func InitServer() {
	cfg := config.GetConfig()
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	v1 := r.Group("/api/v1")
	{
		user := v1.Group("/users")
		routers.User(user)
	}

	r.Run(fmt.Sprintf(":%s", cfg.Server.Port))
}
