package service

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"strings"
	"time"

	"wxcloudrun-golang/internal/model"
	"wxcloudrun-golang/internal/repository"
)

type FamilyDetail struct {
	Family       *model.Family                    `json:"family"`
	CurrentUser  *model.User                      `json:"current_user"`
	CurrentRole  string                           `json:"current_role"`
	Members      []repository.FamilyMemberProfile `json:"members"`
	ActiveInvite *model.FamilyInvite              `json:"active_invite"`
}

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

func GetFamilyDetail(openID string, familyID uint) (*FamilyDetail, error) {
	user, member, err := RequireFamilyMember(openID, familyID)
	if err != nil {
		return nil, err
	}

	family, err := repository.FindFamilyByID(familyID)
	if err != nil {
		return nil, err
	}

	members, err := repository.ListFamilyMembers(familyID)
	if err != nil {
		return nil, err
	}

	invite, err := repository.FindActiveInviteByFamilyID(familyID, time.Now())
	if err != nil && !repository.IsNotFound(err) {
		return nil, err
	}
	if repository.IsNotFound(err) {
		invite = nil
	}

	return &FamilyDetail{
		Family:       family,
		CurrentUser:  user,
		CurrentRole:  member.Role,
		Members:      members,
		ActiveInvite: invite,
	}, nil
}

func CreateFamilyInviteCode(openID string, familyID uint) (*model.FamilyInvite, error) {
	user, member, err := RequireFamilyMember(openID, familyID)
	if err != nil {
		return nil, err
	}
	if member.Role != "owner" {
		return nil, errors.New("仅家庭创建者可以生成邀请码")
	}

	existing, err := repository.FindActiveInviteByFamilyID(familyID, time.Now())
	if err == nil {
		return existing, nil
	}
	if !repository.IsNotFound(err) {
		return nil, err
	}

	code, err := generateInviteCode()
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	invite := &model.FamilyInvite{
		FamilyID:      familyID,
		Code:          code,
		CreatedByUser: user.ID,
		ExpiresAt:     &expiresAt,
	}
	if err := repository.CreateFamilyInvite(invite); err != nil {
		return nil, err
	}
	return invite, nil
}

func UpdateFamilyName(openID string, familyID uint, name string) (*model.Family, error) {
	user, member, err := RequireFamilyMember(openID, familyID)
	if err != nil {
		return nil, err
	}
	if member.Role != "owner" {
		return nil, errors.New("仅家庭创建者可以编辑家庭")
	}

	cleanName := strings.TrimSpace(name)
	if cleanName == "" {
		return nil, errors.New("家庭名称不能为空")
	}

	family, err := repository.FindFamilyByID(familyID)
	if err != nil {
		return nil, err
	}
	if family.OwnerUserID != user.ID {
		return nil, errors.New("仅家庭创建者可以编辑家庭")
	}

	family.Name = cleanName
	if err := repository.SaveFamily(family); err != nil {
		return nil, err
	}
	return family, nil
}

func DeleteFamily(openID string, familyID uint) error {
	user, member, err := RequireFamilyMember(openID, familyID)
	if err != nil {
		return err
	}
	if member.Role != "owner" {
		return errors.New("仅家庭创建者可以删除家庭")
	}

	family, err := repository.FindFamilyByID(familyID)
	if err != nil {
		return err
	}
	if family.OwnerUserID != user.ID {
		return errors.New("仅家庭创建者可以删除家庭")
	}

	return repository.DeleteFamilyByID(familyID)
}

func JoinFamilyByCode(openID string, code string) (*model.Family, error) {
	user, err := CurrentUser(openID)
	if err != nil {
		return nil, err
	}

	cleanCode := strings.ToUpper(strings.TrimSpace(code))
	if cleanCode == "" {
		return nil, errors.New("邀请码不能为空")
	}

	invite, err := repository.FindFamilyInviteByCode(cleanCode)
	if err != nil {
		return nil, errors.New("邀请码不存在")
	}
	if invite.ExpiresAt != nil && !invite.ExpiresAt.After(time.Now()) {
		return nil, errors.New("邀请码已过期")
	}

	member, err := repository.FindFamilyMember(invite.FamilyID, user.ID)
	if err == nil && member != nil {
		family, findErr := repository.FindFamilyByID(invite.FamilyID)
		if findErr != nil {
			return nil, findErr
		}
		return family, nil
	}
	if err != nil && !repository.IsNotFound(err) {
		return nil, err
	}

	newMember := &model.FamilyMember{
		FamilyID: invite.FamilyID,
		UserID:   user.ID,
		Role:     "member",
		Nickname: user.Nickname,
	}
	if err := repository.CreateFamilyMember(newMember); err != nil {
		return nil, err
	}

	invite.UsedCount += 1
	if err := repository.SaveFamilyInvite(invite); err != nil {
		return nil, err
	}

	return repository.FindFamilyByID(invite.FamilyID)
}

func generateInviteCode() (string, error) {
	buf := make([]byte, 5)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	code := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(buf)
	code = strings.ToUpper(strings.ReplaceAll(code, "=", ""))
	if len(code) > 8 {
		code = code[:8]
	}
	return code, nil
}
