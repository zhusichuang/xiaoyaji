package repository

import (
	"wxcloudrun-golang/internal/db"
	"wxcloudrun-golang/internal/model"
)

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
