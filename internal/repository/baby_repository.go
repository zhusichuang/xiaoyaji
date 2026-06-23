package repository

import (
	"wxcloudrun-golang/internal/db"
	"wxcloudrun-golang/internal/model"

	"gorm.io/gorm"
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

func SaveBaby(baby *model.Baby) error {
	return db.Get().Save(baby).Error
}

func DeleteBabyByIDAndFamilyID(babyID, familyID uint) error {
	return db.Get().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("family_id = ? and baby_id = ?", familyID, babyID).Delete(&model.BabyAction{}).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ? and family_id = ?", babyID, familyID).Delete(&model.Baby{}).Error; err != nil {
			return err
		}
		return nil
	})
}
