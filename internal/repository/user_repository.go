package repository

import (
	"wxcloudrun-golang/internal/db"
	"wxcloudrun-golang/internal/model"

	"gorm.io/gorm"
)

func FindUserByOpenID(openID string) (*model.User, error) {
	var user model.User
	err := db.Get().Where("openid = ?", openID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func SaveUser(user *model.User) error {
	return db.Get().Save(user).Error
}

func IsNotFound(err error) bool {
	return err == gorm.ErrRecordNotFound
}
