package handler

import (
	"net/http"
	"strings"

	"wxcloudrun-golang/internal/middleware"
	"wxcloudrun-golang/internal/service"

	"github.com/gin-gonic/gin"
)

type updateUserProfileRequest struct {
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
}

func UpdateCurrentUser(c *gin.Context) {
	var req updateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}

	user, err := service.UpdateCurrentUserProfile(middleware.CurrentOpenID(c), strings.TrimSpace(req.Nickname), strings.TrimSpace(req.AvatarURL))
	if err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}

	success(c, gin.H{"user": user})
}
