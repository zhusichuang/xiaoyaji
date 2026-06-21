package repository

import (
	"wxcloudrun-golang/internal/db"
	"wxcloudrun-golang/internal/model"
)

type FamilyMemberProfile struct {
	ID           uint   `json:"id"`
	FamilyID     uint   `json:"family_id"`
	UserID       uint   `json:"user_id"`
	Role         string `json:"role"`
	Nickname     string `json:"nickname"`
	UserNickname string `json:"user_nickname"`
	AvatarURL    string `json:"avatar_url"`
}

func CreateFamily(family *model.Family) error {
	return db.Get().Create(family).Error
}

func CreateFamilyMember(member *model.FamilyMember) error {
	return db.Get().Create(member).Error
}

func ListFamiliesByUserID(userID uint) ([]model.Family, error) {
	var families []model.Family
	err := db.Get().
		Table("families").
		Joins("join family_members on family_members.family_id = families.id").
		Where("family_members.user_id = ?", userID).
		Order("families.id asc").
		Find(&families).Error
	return families, err
}

func FindFamilyMember(familyID, userID uint) (*model.FamilyMember, error) {
	var member model.FamilyMember
	err := db.Get().Where("family_id = ? and user_id = ?", familyID, userID).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func FindFamilyByID(familyID uint) (*model.Family, error) {
	var family model.Family
	err := db.Get().Where("id = ?", familyID).First(&family).Error
	if err != nil {
		return nil, err
	}
	return &family, nil
}

func ListFamilyMembers(familyID uint) ([]FamilyMemberProfile, error) {
	var members []FamilyMemberProfile
	err := db.Get().
		Table("family_members").
		Select(`
			family_members.id,
			family_members.family_id,
			family_members.user_id,
			family_members.role,
			family_members.nickname,
			users.nickname as user_nickname,
			users.avatar_url
		`).
		Joins("join users on users.id = family_members.user_id").
		Where("family_members.family_id = ?", familyID).
		Order("family_members.id asc").
		Find(&members).Error
	return members, err
}
