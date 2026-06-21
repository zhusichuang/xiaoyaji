package service

import (
	"errors"

	"wxcloudrun-golang/internal/model"
	"wxcloudrun-golang/internal/repository"
)

func CurrentUser(openID string) (*model.User, error) {
	return repository.FindUserByOpenID(openID)
}

func RequireFamilyMember(openID string, familyID uint) (*model.User, *model.FamilyMember, error) {
	user, err := repository.FindUserByOpenID(openID)
	if err != nil {
		return nil, nil, err
	}

	member, err := repository.FindFamilyMember(familyID, user.ID)
	if err != nil {
		return nil, nil, errors.New("无家庭权限")
	}
	return user, member, nil
}
