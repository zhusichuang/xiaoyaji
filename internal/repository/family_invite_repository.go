package repository

import (
	"time"

	"wxcloudrun-golang/internal/db"
	"wxcloudrun-golang/internal/model"
)

func CreateFamilyInvite(invite *model.FamilyInvite) error {
	return db.Get().Create(invite).Error
}

func FindActiveInviteByFamilyID(familyID uint, now time.Time) (*model.FamilyInvite, error) {
	var invite model.FamilyInvite
	err := db.Get().
		Where("family_id = ? and (expires_at is null or expires_at > ?)", familyID, now).
		Order("id desc").
		First(&invite).Error
	if err != nil {
		return nil, err
	}
	return &invite, nil
}

func FindFamilyInviteByCode(code string) (*model.FamilyInvite, error) {
	var invite model.FamilyInvite
	err := db.Get().Where("code = ?", code).First(&invite).Error
	if err != nil {
		return nil, err
	}
	return &invite, nil
}

func SaveFamilyInvite(invite *model.FamilyInvite) error {
	return db.Get().Save(invite).Error
}
