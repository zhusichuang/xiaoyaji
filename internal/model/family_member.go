package model

import "time"

type FamilyMember struct {
	ID        uint      `gorm:"column:id;primaryKey" json:"id"`
	FamilyID  uint      `gorm:"column:family_id;not null;uniqueIndex:idx_family_user" json:"family_id"`
	UserID    uint      `gorm:"column:user_id;not null;uniqueIndex:idx_family_user" json:"user_id"`
	Role      string    `gorm:"column:role;size:32;not null" json:"role"`
	Nickname  string    `gorm:"column:nickname;size:64" json:"nickname"`
	Relation  string    `gorm:"column:relation;size:32" json:"relation"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (FamilyMember) TableName() string {
	return "family_members"
}
