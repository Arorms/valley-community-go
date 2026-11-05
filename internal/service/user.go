package service

import (
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/auth"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/model"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/repository"
)

type UserProfileService interface {
	UpsertFromUser(user auth.User) error
	GetByUserID(userID string) (*model.UserProfile, error)
	GetByUsername(username string) (*model.UserProfile, error)
}

type userProfileService struct {
	profiles repository.UserProfileRepository
}

func NewUserProfileService(profiles repository.UserProfileRepository) UserProfileService {
	return &userProfileService{profiles: profiles}
}

func (s *userProfileService) UpsertFromUser(user auth.User) error {
	return s.profiles.Upsert(&model.UserProfile{
		UserID:    user.Subject,
		Username:  user.Username,
		Nickname:  user.Name,
		AvatarURL: user.Picture,
	})
}

func (s *userProfileService) GetByUserID(userID string) (*model.UserProfile, error) {
	return s.profiles.GetByUserID(userID)
}

func (s *userProfileService) GetByUsername(username string) (*model.UserProfile, error) {
	return s.profiles.GetByUsername(username)
}
