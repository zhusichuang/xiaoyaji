package repository

import (
	"wxcloudrun-golang/internal/db"
	"wxcloudrun-golang/internal/model"
)

func CreateBaby(baby *model.Baby) error {
	return db.Get().Create(baby).Error
}

func ListBabiesByFamilyID(familyID uint) ([]model.Baby, error) {
	var babies []model.Baby
	err := db.Get().Where("family_id = ?", familyID).Order("id asc").Find(&babies).Error
	return babies, err
}

func FindBabyByIDAndFamilyID(babyID, familyID uint) (*model.Baby, error) {
	var baby model.Baby
	err := db.Get().Where("id = ? and family_id = ?", babyID, familyID).First(&baby).Error
	if err != nil {
		return nil, err
	}
	return &baby, nil
}
