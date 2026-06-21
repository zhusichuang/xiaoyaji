package model

import "time"

type Family struct {
	ID          uint      `gorm:"column:id;primaryKey" json:"id"`
	Name        string    `gorm:"column:name;size:128;not null" json:"name"`
	OwnerUserID uint      `gorm:"column:owner_user_id;not null;index" json:"owner_user_id"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (Family) TableName() string {
	return "families"
}
