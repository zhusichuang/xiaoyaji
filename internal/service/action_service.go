package service

import (
	"errors"
	"fmt"
	"time"

	"wxcloudrun-golang/internal/model"
	"wxcloudrun-golang/internal/repository"
	"wxcloudrun-golang/internal/types"
	"wxcloudrun-golang/internal/util"
)

type CreateActionInput struct {
	FamilyID        uint                   `json:"family_id"`
	BabyID          uint                   `json:"baby_id"`
	ActionType      string                 `json:"action_type"`
	ActionTime      string                 `json:"action_time"`
	Summary         string                 `json:"summary"`
	Data            map[string]interface{} `json:"data"`
	Source          string                 `json:"source"`
	ClientRequestID string                 `json:"client_request_id"`
}

type UpdateActionInput struct {
	ActionTime string                 `json:"action_time"`
	Summary    string                 `json:"summary"`
	Data       map[string]interface{} `json:"data"`
}

type ListActionInput struct {
	FamilyID   uint
	BabyID     uint
	ActionType string
	StartTime  string
	EndTime    string
	Limit      int
	Offset     int
}

type ActionView struct {
	ID              uint                   `json:"id"`
	FamilyID        uint                   `json:"family_id"`
	BabyID          uint                   `json:"baby_id"`
	BabyName        string                 `json:"baby_name"`
	ActionType      string                 `json:"action_type"`
	ActionTime      string                 `json:"action_time"`
	Summary         string                 `json:"summary"`
	Data            map[string]interface{} `json:"data"`
	Source          string                 `json:"source"`
	CreatedBy       uint                   `json:"created_by"`
	CreatedByName   string                 `json:"created_by_name"`
	ClientRequestID string                 `json:"client_request_id"`
	CreatedAt       string                 `json:"created_at"`
	UpdatedAt       string                 `json:"updated_at"`
}

func CreateAction(openID string, input CreateActionInput) (*model.BabyAction, error) {
	user, _, err := RequireFamilyMember(openID, input.FamilyID)
	if err != nil {
		return nil, err
	}

	baby, err := repository.FindBabyByIDAndFamilyID(input.BabyID, input.FamilyID)
	if err != nil {
		return nil, errors.New("宝宝不存在")
	}

	if input.ClientRequestID != "" {
		existed, err := repository.FindActionByRequestID(input.FamilyID, input.ClientRequestID)
		if err == nil {
			return existed, nil
		}
		if !repository.IsNotFound(err) {
			return nil, err
		}
	}

	actionTime, err := util.ParseRFC3339(input.ActionTime)
	if err != nil {
		return nil, fmt.Errorf("action_time 格式错误: %w", err)
	}

	action := &model.BabyAction{
		FamilyID:        input.FamilyID,
		BabyID:          input.BabyID,
		BabyName:        displayBabyName(baby),
		ActionType:      input.ActionType,
		ActionTime:      actionTime,
		Summary:         input.Summary,
		DataJSON:        util.MustMarshal(input.Data),
		Source:          defaultSource(input.Source),
		CreatedBy:       user.ID,
		CreatedByName:   user.Nickname,
		ClientRequestID: input.ClientRequestID,
	}
	if err := repository.CreateAction(action); err != nil {
		return nil, err
	}
	return action, nil
}

func BatchCreateActions(openID string, familyID uint, records []types.RecordPayload) ([]model.BabyAction, error) {
	result := make([]model.BabyAction, 0, len(records))
	for _, record := range records {
		action, err := CreateAction(openID, CreateActionInput{
			FamilyID:        familyID,
			BabyID:          record.BabyID,
			ActionType:      record.ActionType,
			ActionTime:      record.ActionTime,
			Summary:         record.Summary,
			Data:            record.Data,
			Source:          record.Source,
			ClientRequestID: record.ClientRequestID,
		})
		if err != nil {
			return nil, err
		}
		result = append(result, *action)
	}
	return result, nil
}

func GetAction(openID string, actionID uint) (ActionView, error) {
	action, err := repository.FindActionByID(actionID)
	if err != nil {
		return ActionView{}, err
	}

	if _, _, err := RequireFamilyMember(openID, action.FamilyID); err != nil {
		return ActionView{}, err
	}

	return buildActionView(*action), nil
}

func UpdateAction(openID string, actionID uint, input UpdateActionInput) (ActionView, error) {
	action, err := repository.FindActionByID(actionID)
	if err != nil {
		return ActionView{}, err
	}

	if _, _, err := RequireFamilyMember(openID, action.FamilyID); err != nil {
		return ActionView{}, err
	}

	actionTime, err := util.ParseRFC3339(input.ActionTime)
	if err != nil {
		return ActionView{}, fmt.Errorf("action_time 格式错误: %w", err)
	}

	action.ActionTime = actionTime
	action.Summary = input.Summary
	action.DataJSON = util.MustMarshal(input.Data)
	if err := repository.UpdateAction(action); err != nil {
		return ActionView{}, err
	}

	return buildActionView(*action), nil
}

func DeleteAction(openID string, actionID uint) error {
	action, err := repository.FindActionByID(actionID)
	if err != nil {
		return err
	}

	if _, _, err := RequireFamilyMember(openID, action.FamilyID); err != nil {
		return err
	}

	return repository.SoftDeleteAction(actionID)
}

func ListActions(openID string, input ListActionInput) ([]ActionView, error) {
	if _, _, err := RequireFamilyMember(openID, input.FamilyID); err != nil {
		return nil, err
	}

	var startPtr *time.Time
	if input.StartTime != "" {
		start, err := util.ParseRFC3339(input.StartTime)
		if err != nil {
			return nil, err
		}
		startPtr = &start
	}

	var endPtr *time.Time
	if input.EndTime != "" {
		end, err := util.ParseRFC3339(input.EndTime)
		if err != nil {
			return nil, err
		}
		endPtr = &end
	}

	if input.Limit <= 0 {
		input.Limit = 20
	}
	if input.Limit > 100 {
		input.Limit = 100
	}

	actions, err := repository.ListActions(repository.ListActionsInput{
		FamilyID:   input.FamilyID,
		BabyID:     input.BabyID,
		ActionType: input.ActionType,
		StartTime:  startPtr,
		EndTime:    endPtr,
		Limit:      input.Limit,
		Offset:     input.Offset,
	})
	if err != nil {
		return nil, err
	}

	views := make([]ActionView, 0, len(actions))
	for _, action := range actions {
		views = append(views, buildActionView(action))
	}
	return views, nil
}

func GetTodaySummary(openID string, familyID, babyID uint, timezoneOffsetMin int) (map[string]interface{}, error) {
	if _, _, err := RequireFamilyMember(openID, familyID); err != nil {
		return nil, err
	}

	start, end, dateText := util.TodayRangeByOffset(timezoneOffsetMin)
	actions, err := repository.ListActions(repository.ListActionsInput{
		FamilyID:  familyID,
		BabyID:    babyID,
		StartTime: &start,
		EndTime:   &end,
		Limit:     500,
		Offset:    0,
	})
	if err != nil {
		return nil, err
	}

	summary := map[string]interface{}{
		"baby_id":         babyID,
		"date":            dateText,
		"feed_total_ml":   0,
		"feed_count":      0,
		"pee_count":       0,
		"poop_count":      0,
		"sleep_total_min": 0,
		"last_action":     nil,
	}

	if len(actions) > 0 {
		summary["last_action"] = buildActionView(actions[0])
	}

	for _, action := range actions {
		data := util.MustUnmarshal(action.DataJSON)
		switch action.ActionType {
		case "feed":
			summary["feed_count"] = summary["feed_count"].(int) + 1
			summary["feed_total_ml"] = summary["feed_total_ml"].(int) + toInt(data["amount_ml"])
		case "diaper":
			if toBool(data["pee"]) {
				summary["pee_count"] = summary["pee_count"].(int) + 1
			}
			if toBool(data["poop"]) {
				summary["poop_count"] = summary["poop_count"].(int) + 1
			}
		case "sleep":
			summary["sleep_total_min"] = summary["sleep_total_min"].(int) + toInt(data["duration_min"])
		}
	}

	return summary, nil
}

func buildActionView(action model.BabyAction) ActionView {
	return ActionView{
		ID:              action.ID,
		FamilyID:        action.FamilyID,
		BabyID:          action.BabyID,
		BabyName:        action.BabyName,
		ActionType:      action.ActionType,
		ActionTime:      action.ActionTime.UTC().Format(time.RFC3339),
		Summary:         action.Summary,
		Data:            util.MustUnmarshal(action.DataJSON),
		Source:          action.Source,
		CreatedBy:       action.CreatedBy,
		CreatedByName:   action.CreatedByName,
		ClientRequestID: action.ClientRequestID,
		CreatedAt:       action.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:       action.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func defaultSource(source string) string {
	if source == "" {
		return "manual"
	}
	return source
}

func displayBabyName(baby *model.Baby) string {
	if baby.Nickname != "" {
		return baby.Nickname
	}
	return baby.Name
}

func toInt(value interface{}) int {
	switch typed := value.(type) {
	case float64:
		return int(typed)
	case float32:
		return int(typed)
	case int:
		return typed
	case int64:
		return int(typed)
	case string:
		var parsed int
		fmt.Sscanf(typed, "%d", &parsed)
		return parsed
	default:
		return 0
	}
}

func toBool(value interface{}) bool {
	typed, ok := value.(bool)
	return ok && typed
}
