package model

import "time"

type FamilyInvite struct {
	ID            uint       `gorm:"column:id;primaryKey" json:"id"`
	FamilyID      uint       `gorm:"column:family_id;not null;index" json:"family_id"`
	Code          string     `gorm:"column:code;size:16;not null;uniqueIndex" json:"code"`
	CreatedByUser uint       `gorm:"column:created_by_user;not null" json:"created_by_user"`
	ExpiresAt     *time.Time `gorm:"column:expires_at" json:"expires_at"`
	UsedCount     uint       `gorm:"column:used_count;not null;default:0" json:"used_count"`
	CreatedAt     time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (FamilyInvite) TableName() string {
	return "family_invites"
}
