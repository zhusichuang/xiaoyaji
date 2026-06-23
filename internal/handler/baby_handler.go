package handler

import (
	"errors"
	"net/http"
	"strconv"

	"wxcloudrun-golang/internal/middleware"
	"wxcloudrun-golang/internal/service"

	"github.com/gin-gonic/gin"
)

func CreateBaby(c *gin.Context) {
	var req service.CreateBabyInput
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}

	baby, err := service.CreateBaby(middleware.CurrentOpenID(c), req)
	if err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}
	success(c, gin.H{"baby_id": baby.ID})
}

func ListBabies(c *gin.Context) {
	familyID, err := strconv.ParseUint(c.Query("family_id"), 10, 64)
	if err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}

	babies, err := service.ListBabies(middleware.CurrentOpenID(c), uint(familyID))
	if err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}
	success(c, gin.H{"babies": babies})
}

func UpdateBaby(c *gin.Context) {
	babyID, err := strconv.ParseUint(c.Param("babyID"), 10, 64)
	if err != nil || babyID == 0 {
		fail(c, http.StatusBadRequest, errors.New("baby_id 非法"))
		return
	}

	var req service.CreateBabyInput
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}
	req.ID = uint(babyID)

	baby, err := service.UpdateBaby(middleware.CurrentOpenID(c), req)
	if err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}
	success(c, gin.H{"baby": baby})
}

func DeleteBaby(c *gin.Context) {
	familyID, err := strconv.ParseUint(c.Query("family_id"), 10, 64)
	if err != nil || familyID == 0 {
		fail(c, http.StatusBadRequest, errors.New("family_id 非法"))
		return
	}

	babyID, err := strconv.ParseUint(c.Param("babyID"), 10, 64)
	if err != nil || babyID == 0 {
		fail(c, http.StatusBadRequest, errors.New("baby_id 非法"))
		return
	}

	if err := service.DeleteBaby(middleware.CurrentOpenID(c), uint(familyID), uint(babyID)); err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}
	success(c, gin.H{"deleted": true})
}
