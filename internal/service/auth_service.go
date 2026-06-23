package service

import (
	"time"

	"wxcloudrun-golang/internal/model"
	"wxcloudrun-golang/internal/repository"
)

func EnsureUserAndDefaultFamily(openID string) (*model.User, []model.Family, error) {
	user, err := repository.FindUserByOpenID(openID)
	if err != nil {
		if !repository.IsNotFound(err) {
			return nil, nil, err
		}

		user = &model.User{
			OpenID:      openID,
			Nickname:    "微信用户",
			LastLoginAt: time.Now(),
		}
		if err := repository.SaveUser(user); err != nil {
			return nil, nil, err
		}
	} else {
		user.LastLoginAt = time.Now()
		if err := repository.SaveUser(user); err != nil {
			return nil, nil, err
		}
	}

	families, err := repository.ListFamiliesByUserID(user.ID)
	if err != nil {
		return nil, nil, err
	}
	if len(families) > 0 {
		return user, families, nil
	}

	family := &model.Family{
		Name:        "我的家庭",
		OwnerUserID: user.ID,
	}
	if err := repository.CreateFamily(family); err != nil {
		return nil, nil, err
	}

	member := &model.FamilyMember{
		FamilyID: family.ID,
		UserID:   user.ID,
		Role:     "owner",
		Nickname: user.Nickname,
	}
	if err := repository.CreateFamilyMember(member); err != nil {
		return nil, nil, err
	}

	return user, []model.Family{*family}, nil
}
