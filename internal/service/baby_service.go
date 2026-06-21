package service

import (
	"wxcloudrun-golang/internal/model"
	"wxcloudrun-golang/internal/repository"
)

type CreateBabyInput struct {
	FamilyID        uint   `json:"family_id"`
	Name            string `json:"name"`
	Nickname        string `json:"nickname"`
	Gender          string `json:"gender"`
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
