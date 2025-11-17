package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/mohsen104/web-api/api/handlers"
	"github.com/mohsen104/web-api/api/middlewares"
	"github.com/mohsen104/web-api/config"
)

func User(r *gin.RouterGroup, cfg *config.Config) {
	handler := handlers.NewUserHandler(cfg)

	r.POST("/send-otp", middlewares.OtpLimiter(cfg), handler.SendOtp)
}
