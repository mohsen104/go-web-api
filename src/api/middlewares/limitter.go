package middlewares

import (
	"net/http"

	"github.com/didip/tollbooth"
	"github.com/gin-gonic/gin"
	"github.com/mohsen104/web-api/api/helper"
)

func LimitByRequest() gin.HandlerFunc {
	lmt := tollbooth.NewLimiter(5, nil)
	return func(ctx *gin.Context) {
		err := tollbooth.LimitByRequest(lmt, ctx.Writer, ctx.Request)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, helper.GenerateBaseResponseWithError("Too many requests", false, 0, err))
			return
		}
		ctx.Next()
	}
}
