package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mohsen104/web-api/api/dto"
	"github.com/mohsen104/web-api/api/helper"
	"github.com/mohsen104/web-api/config"
	"github.com/mohsen104/web-api/services"
)

type CountryHandler struct {
	service *services.CountryService
}

func NewCountryHandler(cfg *config.Config) *CountryHandler {
	return &CountryHandler{service: services.NewCountryService(cfg)}
}

// CreateCountry
// @Summary Create a new country
// @Description Create a new country with the provided details
// @Tags countries
// @Accept json
// @Produce json
// @Param country body dto.CreateUpdateCountryRequest true "Country object to be created"
// @Success 201 {object} dto.CountryResponse
// @Failure 400 {object} helper.BaseResponseWithError
// @Failure 500 {object} helper.BaseResponseWithError
// @Router /countries [post]
func (h *CountryHandler) Create(c *gin.Context) {
	req := dto.CreateUpdateCountryRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, -1, err),
		)
		return
	}

	res, err := h.service.Create(c, &req)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, -1, err),
		)
		return
	}
	c.JSON(http.StatusCreated,
		helper.GenerateBaseResponse(res, true, 0),
	)
}

// UpdateCountry
// @Summary Update a country
// @Description Update a country with the provided details
// @Tags countries
// @Accept json
// @Produce json
// @Param country body dto.CreateUpdateCountryRequest true "Country object to be created"
// @Success 201 {object} dto.CountryResponse
// @Failure 400 {object} helper.BaseResponseWithError
// @Failure 500 {object} helper.BaseResponseWithError
// @Router /countries [post]
func (h *CountryHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	req := dto.CreateUpdateCountryRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, -1, err),
		)
		return
	}

	res, err := h.service.Update(c, id, &req)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, -1, err),
		)
		return
	}
	c.JSON(http.StatusOK,
		helper.GenerateBaseResponse(res, true, 0),
	)
}
