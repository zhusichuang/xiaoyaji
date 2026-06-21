package model

import "time"

type Baby struct {
	ID              uint      `gorm:"column:id;primaryKey" json:"id"`
	FamilyID        uint      `gorm:"column:family_id;not null;index" json:"family_id"`
	Name            string    `gorm:"column:name;size:128;not null" json:"name"`
	Nickname        string    `gorm:"column:nickname;size:128" json:"nickname"`
	Gender          string    `gorm:"column:gender;size:32" json:"gender"`
	BirthDate       string    `gorm:"column:birth_date;size:32" json:"birth_date"`
	BirthTime       string    `gorm:"column:birth_time;size:32" json:"birth_time"`
	GestationalWeek int       `gorm:"column:gestational_week" json:"gestational_week"`
	BirthWeightG    int       `gorm:"column:birth_weight_g" json:"birth_weight_g"`
	AvatarURL       string    `gorm:"column:avatar_url;size:255" json:"avatar_url"`
	Remark          string    `gorm:"column:remark;type:text" json:"remark"`
	CreatedBy       uint      `gorm:"column:created_by;not null" json:"created_by"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (Baby) TableName() string {
	return "babies"
}
