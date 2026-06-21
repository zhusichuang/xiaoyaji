package handler

import (
	"errors"
	"net/http"

	"wxcloudrun-golang/internal/middleware"
	"wxcloudrun-golang/internal/model"
	"wxcloudrun-golang/internal/repository"
	"wxcloudrun-golang/internal/service"

	"github.com/gin-gonic/gin"
)

type createFamilyRequest struct {
	Name string `json:"name"`
}

func ListFamilies(c *gin.Context) {
	user, err := service.CurrentUser(middleware.CurrentOpenID(c))
	if err != nil {
		fail(c, http.StatusInternalServerError, err)
		return
	}

	families, err := repository.ListFamiliesByUserID(user.ID)
	if err != nil {
		fail(c, http.StatusInternalServerError, err)
		return
	}
	success(c, gin.H{"families": families})
}

func CreateFamily(c *gin.Context) {
	var req createFamilyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}
	if req.Name == "" {
		fail(c, http.StatusBadRequest, errors.New("name 不能为空"))
		return
	}

	user, err := service.CurrentUser(middleware.CurrentOpenID(c))
	if err != nil {
		fail(c, http.StatusInternalServerError, err)
		return
	}

	family := &model.Family{
		Name:        req.Name,
		OwnerUserID: user.ID,
	}
	if err := repository.CreateFamily(family); err != nil {
		fail(c, http.StatusInternalServerError, err)
		return
	}

	member := &model.FamilyMember{
		FamilyID: family.ID,
		UserID:   user.ID,
		Role:     "owner",
		Nickname: "我",
	}
	if err := repository.CreateFamilyMember(member); err != nil {
		fail(c, http.StatusInternalServerError, err)
		return
	}

	success(c, gin.H{"family_id": family.ID})
}
