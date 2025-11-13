package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TestHandler struct{}

type header struct{}

func NewTestHandler() *TestHandler {
	return &TestHandler{}
}

func (t *TestHandler) Test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"result": "Test",
	})
}

func (t *TestHandler) TestId(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"result": "Test",
		"Id":     id,
	})
}

func (t *TestHandler) Header(c *gin.Context) {
	// userId := c.GetHeader("UserId")
	header := header{}
	c.BindHeader(&header)
	c.JSON(http.StatusOK, gin.H{
		"result": "Test",
		"Header": header,
	})
}

func (t *TestHandler) Query(c *gin.Context) {
	name := c.BindQuery("name")

	c.JSON(http.StatusOK, gin.H{
		"result": "Test",
		"name":   name,
	})
}

func (t *TestHandler) Create(c *gin.Context) {
	p := struct {
		Name   string `json:"name" binding:"required,min=3,max=10"`
		Mobile string `json:"mobile" binding:"required"`
	}{}

	err := c.ShouldBindJSON(&p)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "Test",
		"name":   p.Name,
	})
}
