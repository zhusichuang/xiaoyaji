package handler

import (
	"errors"
	"net/http"
	"strconv"

	"wxcloudrun-golang/internal/middleware"
	"wxcloudrun-golang/internal/model"
	"wxcloudrun-golang/internal/repository"
	"wxcloudrun-golang/internal/service"

	"github.com/gin-gonic/gin"
)

type createFamilyRequest struct {
	Name string `json:"name"`
}

type joinFamilyRequest struct {
	Code string `json:"code"`
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

func UpdateFamily(c *gin.Context) {
	familyID, err := strconv.ParseUint(c.Param("familyID"), 10, 64)
	if err != nil || familyID == 0 {
		fail(c, http.StatusBadRequest, errors.New("family_id 非法"))
		return
	}

	var req createFamilyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}

	family, err := service.UpdateFamilyName(middleware.CurrentOpenID(c), uint(familyID), req.Name)
	if err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}
	success(c, gin.H{"family": family})
}

func DeleteFamily(c *gin.Context) {
	familyID, err := strconv.ParseUint(c.Param("familyID"), 10, 64)
	if err != nil || familyID == 0 {
		fail(c, http.StatusBadRequest, errors.New("family_id 非法"))
		return
	}

	if err := service.DeleteFamily(middleware.CurrentOpenID(c), uint(familyID)); err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}
	success(c, gin.H{"deleted": true})
}

func GetFamilyDetail(c *gin.Context) {
	familyID, err := strconv.ParseUint(c.Param("familyID"), 10, 64)
	if err != nil || familyID == 0 {
		fail(c, http.StatusBadRequest, errors.New("family_id 非法"))
		return
	}

	detail, err := service.GetFamilyDetail(middleware.CurrentOpenID(c), uint(familyID))
	if err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}

	success(c, detail)
}

func CreateFamilyInviteCode(c *gin.Context) {
	familyID, err := strconv.ParseUint(c.Param("familyID"), 10, 64)
	if err != nil || familyID == 0 {
		fail(c, http.StatusBadRequest, errors.New("family_id 非法"))
		return
	}

	invite, err := service.CreateFamilyInviteCode(middleware.CurrentOpenID(c), uint(familyID))
	if err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}

	success(c, gin.H{"invite": invite})
}

func JoinFamily(c *gin.Context) {
	var req joinFamilyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}

	family, err := service.JoinFamilyByCode(middleware.CurrentOpenID(c), req.Code)
	if err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}

	success(c, gin.H{"family": family})
}
