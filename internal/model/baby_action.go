package model

import "time"

type BabyAction struct {
	ID              uint      `gorm:"column:id;primaryKey" json:"id"`
	FamilyID        uint      `gorm:"column:family_id;not null;index:idx_family_baby_time,priority:1" json:"family_id"`
	BabyID          uint      `gorm:"column:baby_id;not null;index:idx_family_baby_time,priority:2" json:"baby_id"`
	BabyName        string    `gorm:"column:baby_name;size:128" json:"baby_name"`
	ActionType      string    `gorm:"column:action_type;size:32;not null;index" json:"action_type"`
	ActionTime      time.Time `gorm:"column:action_time;not null;index:idx_family_baby_time,priority:3" json:"action_time"`
	Summary         string    `gorm:"column:summary;size:255;not null" json:"summary"`
	DataJSON        string    `gorm:"column:data_json;type:text" json:"-"`
	Source          string    `gorm:"column:source;size:32;not null" json:"source"`
	CreatedBy       uint      `gorm:"column:created_by;not null" json:"created_by"`
	CreatedByName   string    `gorm:"column:created_by_name;size:128" json:"created_by_name"`
	ClientRequestID string    `gorm:"column:client_request_id;size:128;index" json:"client_request_id"`
	Deleted         bool      `gorm:"column:deleted;default:false;index" json:"deleted"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (BabyAction) TableName() string {
	return "baby_actions"
}
