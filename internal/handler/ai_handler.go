package handler

import (
	"net/http"

	"wxcloudrun-golang/internal/middleware"
	"wxcloudrun-golang/internal/service"

	"github.com/gin-gonic/gin"
)

func ParseRecord(c *gin.Context) {
	var req service.ParseRecordInput
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}

	result, err := service.ParseRecord(middleware.CurrentOpenID(c), req)
	if err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}
	success(c, result)
}

func Chat(c *gin.Context) {
	var req service.ChatInput
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}

	result, err := service.Chat(middleware.CurrentOpenID(c), req)
	if err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}
	success(c, result)
}
