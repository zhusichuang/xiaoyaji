package handler

import (
	"net/http"
	"strconv"

	"wxcloudrun-golang/internal/middleware"
	"wxcloudrun-golang/internal/service"
	"wxcloudrun-golang/internal/types"

	"github.com/gin-gonic/gin"
)

type batchCreateRequest struct {
	FamilyID uint                  `json:"family_id"`
	Records  []types.RecordPayload `json:"records"`
}

func CreateAction(c *gin.Context) {
	var req service.CreateActionInput
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}

	action, err := service.CreateAction(middleware.CurrentOpenID(c), req)
	if err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}
	success(c, gin.H{"action_id": action.ID})
}

func GetAction(c *gin.Context) {
	actionID, err := strconv.ParseUint(c.Param("actionID"), 10, 64)
	if err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}

	action, err := service.GetAction(middleware.CurrentOpenID(c), uint(actionID))
	if err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}
	success(c, action)
}

func UpdateAction(c *gin.Context) {
	actionID, err := strconv.ParseUint(c.Param("actionID"), 10, 64)
	if err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}

	var req service.UpdateActionInput
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}

	action, err := service.UpdateAction(middleware.CurrentOpenID(c), uint(actionID), req)
	if err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}
	success(c, gin.H{"action": action})
}

func DeleteAction(c *gin.Context) {
	actionID, err := strconv.ParseUint(c.Param("actionID"), 10, 64)
	if err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}

	if err := service.DeleteAction(middleware.CurrentOpenID(c), uint(actionID)); err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}
	success(c, gin.H{"deleted": true})
}

func BatchCreateActions(c *gin.Context) {
	var req batchCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}

	actions, err := service.BatchCreateActions(middleware.CurrentOpenID(c), req.FamilyID, req.Records)
	if err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}

	ids := make([]uint, 0, len(actions))
	for _, action := range actions {
		ids = append(ids, action.ID)
	}
	success(c, gin.H{"created_count": len(actions), "action_ids": ids})
}

func ListActions(c *gin.Context) {
	familyID, err := strconv.ParseUint(c.Query("family_id"), 10, 64)
	if err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}
	babyID, _ := strconv.ParseUint(c.Query("baby_id"), 10, 64)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	actions, err := service.ListActions(middleware.CurrentOpenID(c), service.ListActionInput{
		FamilyID:   uint(familyID),
		BabyID:     uint(babyID),
		ActionType: c.Query("action_type"),
		StartTime:  c.Query("start_time"),
		EndTime:    c.Query("end_time"),
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}
	success(c, gin.H{"actions": actions, "has_more": len(actions) == limit})
}

func TodaySummary(c *gin.Context) {
	familyID, err := strconv.ParseUint(c.Query("family_id"), 10, 64)
	if err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}
	babyID, err := strconv.ParseUint(c.Query("baby_id"), 10, 64)
	if err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}
	offset, _ := strconv.Atoi(c.DefaultQuery("timezone_offset_min", "0"))

	summary, err := service.GetTodaySummary(middleware.CurrentOpenID(c), uint(familyID), uint(babyID), offset)
	if err != nil {
		fail(c, http.StatusBadRequest, err)
		return
	}
	success(c, summary)
}
