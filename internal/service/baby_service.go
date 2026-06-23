package service

import (
	"errors"

	"wxcloudrun-golang/internal/model"
	"wxcloudrun-golang/internal/repository"
)

type CreateBabyInput struct {
	FamilyID        uint   `json:"family_id"`
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	Nickname        string `json:"nickname"`
	Gender          string `json:"gender"`
	AvatarURL       string `json:"avatar_url"`
	BirthDate       string `json:"birth_date"`
	BirthTime       string `json:"birth_time"`
	GestationalWeek int    `json:"gestational_week"`
	BirthWeightG    int    `json:"birth_weight_g"`
	Remark          string `json:"remark"`
}

func CreateBaby(openID string, input CreateBabyInput) (*model.Baby, error) {
	user, _, err := RequireFamilyMember(openID, input.FamilyID)
	if err != nil {
		return nil, err
	}

	baby := &model.Baby{
		FamilyID:        input.FamilyID,
		Name:            input.Name,
		Nickname:        input.Nickname,
		Gender:          input.Gender,
		AvatarURL:       input.AvatarURL,
		BirthDate:       input.BirthDate,
		BirthTime:       input.BirthTime,
		GestationalWeek: input.GestationalWeek,
		BirthWeightG:    input.BirthWeightG,
		Remark:          input.Remark,
		CreatedBy:       user.ID,
	}
	if err := repository.CreateBaby(baby); err != nil {
		return nil, err
	}
	return baby, nil
}

func ListBabies(openID string, familyID uint) ([]model.Baby, error) {
	if _, _, err := RequireFamilyMember(openID, familyID); err != nil {
		return nil, err
	}
	return repository.ListBabiesByFamilyID(familyID)
}

func UpdateBaby(openID string, input CreateBabyInput) (*model.Baby, error) {
	if input.ID == 0 {
		return nil, errors.New("baby_id 不能为空")
	}
	if input.FamilyID == 0 {
		return nil, errors.New("family_id 不能为空")
	}

	if _, _, err := RequireFamilyMember(openID, input.FamilyID); err != nil {
		return nil, err
	}

	baby, err := repository.FindBabyByIDAndFamilyID(input.ID, input.FamilyID)
	if err != nil {
		return nil, err
	}

	baby.Name = input.Name
	baby.Nickname = input.Nickname
	baby.Gender = input.Gender
	baby.AvatarURL = input.AvatarURL
	baby.BirthDate = input.BirthDate
	baby.BirthTime = input.BirthTime
	baby.GestationalWeek = input.GestationalWeek
	baby.BirthWeightG = input.BirthWeightG
	baby.Remark = input.Remark

	if err := repository.SaveBaby(baby); err != nil {
		return nil, err
	}
	return baby, nil
}

func DeleteBaby(openID string, familyID, babyID uint) error {
	if _, _, err := RequireFamilyMember(openID, familyID); err != nil {
		return err
	}
	return repository.DeleteBabyByIDAndFamilyID(babyID, familyID)
}
