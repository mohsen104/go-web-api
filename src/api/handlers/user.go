package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohsen104/web-api/api/helper"
	"github.com/mohsen104/web-api/config"
	"github.com/mohsen104/web-api/services"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(cfg *config.Config) *UserHandler {
	service := services.NewUserService(cfg)
	return &UserHandler{service: service}
}

// SendOtp godoc
// @Summary Send otp
// @Description Send otp
// @Tags users
// @Accept json
// @Produce json
// @Param request body services.GetOtpRequest true "GetOtpRequest"
// @Success 201 {object} helper.BaseHttpResponse "Success"
// @Failure 400 {object} helper.BaseHttpResponse "Failed"
// @Failure 409 {object} helper.BaseHttpResponse "Failed"
// @Router /v1/users/send-otp [post]
func (h *UserHandler) SendOtp(c *gin.Context) {
	otp := new(services.GetOtpRequest)
	err := c.ShouldBindJSON(&otp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, -1, err),
		)
		return
	}
	err = h.service.SendOtp(otp)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, -1, err),
		)
		return
	}
	c.JSON(http.StatusCreated,
		helper.GenerateBaseResponse(nil, true, 0),
	)
}
