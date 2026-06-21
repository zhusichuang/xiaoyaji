package handler

import (
	"net/http"

	"wxcloudrun-golang/internal/middleware"
	"wxcloudrun-golang/internal/model"
	"wxcloudrun-golang/internal/repository"
	"wxcloudrun-golang/internal/service"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	openID := middleware.CurrentOpenID(c)
	user, families, err := service.EnsureUserAndDefaultFamily(openID)
	if err != nil {
		fail(c, http.StatusInternalServerError, err)
		return
	}

	var babies interface{} = []interface{}{}
	var currentBabyID uint
	if len(families) > 0 {
		list, err := repository.ListBabiesByFamilyID(families[0].ID)
		if err != nil {
			fail(c, http.StatusInternalServerError, err)
			return
		}
		babies = list
		if len(list) > 0 {
			currentBabyID = list[0].ID
		}
	}

	success(c, gin.H{
		"user":            user,
		"current_family":  firstFamily(families),
		"families":        families,
		"babies":          babies,
		"current_baby_id": currentBabyID,
	})
}

func firstFamily(families []model.Family) interface{} {
	if len(families) == 0 {
		return nil
	}
	return families[0]
}
