package repository

import (
	"time"

	"wxcloudrun-golang/internal/db"
	"wxcloudrun-golang/internal/model"
)

type ListActionsInput struct {
	FamilyID   uint
	BabyID     uint
	ActionType string
	StartTime  *time.Time
	EndTime    *time.Time
	Limit      int
	Offset     int
}

func CreateAction(action *model.BabyAction) error {
	return db.Get().Create(action).Error
}

func FindActionByRequestID(familyID uint, requestID string) (*model.BabyAction, error) {
	var action model.BabyAction
	err := db.Get().Where("family_id = ? and client_request_id = ? and deleted = ?", familyID, requestID, false).First(&action).Error
	if err != nil {
		return nil, err
	}
	return &action, nil
}

func FindActionByID(actionID uint) (*model.BabyAction, error) {
	var action model.BabyAction
	err := db.Get().Where("id = ? and deleted = ?", actionID, false).First(&action).Error
	if err != nil {
		return nil, err
	}
	return &action, nil
}

func UpdateAction(action *model.BabyAction) error {
	return db.Get().Save(action).Error
}

func SoftDeleteAction(actionID uint) error {
	return db.Get().Model(&model.BabyAction{}).Where("id = ?", actionID).Updates(map[string]interface{}{
		"deleted": true,
	}).Error
}

func ListActions(input ListActionsInput) ([]model.BabyAction, error) {
	query := db.Get().Where("family_id = ? and deleted = ?", input.FamilyID, false)
	if input.BabyID > 0 {
		query = query.Where("baby_id = ?", input.BabyID)
	}
	if input.ActionType != "" {
		query = query.Where("action_type = ?", input.ActionType)
	}
	if input.StartTime != nil {
		query = query.Where("action_time >= ?", *input.StartTime)
	}
	if input.EndTime != nil {
		query = query.Where("action_time <= ?", *input.EndTime)
	}

	var actions []model.BabyAction
	err := query.Order("action_time desc").Limit(input.Limit).Offset(input.Offset).Find(&actions).Error
	return actions, err
}
