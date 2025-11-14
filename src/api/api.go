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
	"github.com/mohsen104/web-api/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitServer(cfg *config.Config) {
	r := gin.New()

	val, ok := binding.Validator.Engine().(*validator.Validate)

	if ok {
		val.RegisterValidation("mobile", validations.IranianMobileNumberValidator)
	}

	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(middlewares.Cors(cfg))
	r.Use(middlewares.LimitByRequest())

	RegisterSwagger(r, cfg)

	RegisterRouters(r)

	r.Run(fmt.Sprintf(":%s", cfg.Server.Port))
}

func RegisterSwagger(r *gin.Engine, cfg *config.Config) {
	docs.SwaggerInfo.Title = "Go Web API"
	docs.SwaggerInfo.Description = "Go Web API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.Server.Domain, cfg.Server.Port)
	docs.SwaggerInfo.Schemes = []string{"http"}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func RegisterRouters(r *gin.Engine) {
	api := r.Group("/api")
	v1 := api.Group("/v1")

	{
		user := v1.Group("/users")
		routers.User(user)
	}
}
