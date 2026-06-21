package model

import "time"

type User struct {
	ID          uint      `gorm:"column:id;primaryKey" json:"id"`
	OpenID      string    `gorm:"column:openid;size:128;uniqueIndex;not null" json:"openid"`
	UnionID     string    `gorm:"column:unionid;size:128" json:"unionid"`
	Nickname    string    `gorm:"column:nickname;size:128;not null" json:"nickname"`
	AvatarURL   string    `gorm:"column:avatar_url;size:255" json:"avatar_url"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	LastLoginAt time.Time `gorm:"column:last_login_at" json:"last_login_at"`
}

func (User) TableName() string {
	return "users"
}
