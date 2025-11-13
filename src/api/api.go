package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/mohsen104/web-api/api/middlewares"
	"github.com/mohsen104/web-api/api/routers"
	"github.com/mohsen104/web-api/api/validations"
	"github.com/mohsen104/web-api/config"
)

func InitServer() {
	cfg := config.GetConfig()
	r := gin.New()
	val, ok := binding.Validator.Engine().(*validator.Validate)

	if ok {
		val.RegisterValidation("mobile", validations.IranianMobileNumberValidator)
	}

	r.Use(gin.Logger(), gin.Recovery(), middlewares.LimitByRequest())

	api := r.Group("/api")
	v1 := api.Group("/v1")
	{
		user := v1.Group("/users")
		routers.User(user)

		test := v1.Group("/test")
		routers.TestRouter(test)
	}

	r.Run(fmt.Sprintf(":%s", cfg.Server.Port))
}
