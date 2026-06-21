package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": data,
	})
}

func fail(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{
		"code":     -1,
		"errorMsg": err.Error(),
		"data":     nil,
	})
}
