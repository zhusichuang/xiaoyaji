package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	content, err := os.ReadFile("./index.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "internal error")
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", content)
}
