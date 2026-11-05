package repository

import (
	"errors"

	"gorm.io/gorm"

	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/model"
)

type UserProfileRepository interface {
	Upsert(profile *model.UserProfile) error
	GetByUserID(userID string) (*model.UserProfile, error)
	GetByUsername(username string) (*model.UserProfile, error)
	ListByIDs(userIDs []string) (map[string]model.UserProfile, error)
}

type userProfileRepository struct {
	db *gorm.DB
}

func NewUserProfileRepository(db *gorm.DB) UserProfileRepository {
	return &userProfileRepository{db: db}
}

// Upsert 同步 Keycloak 资料；username 冲突时保留旧用户名，不阻断登录。
func (r *userProfileRepository) Upsert(profile *model.UserProfile) error {
	var existing model.UserProfile
	err := r.db.Where("user_id = ?", profile.UserID).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.db.Create(profile).Error
	}
	if err != nil {
		return err
	}

	var conflict int64
	if err := r.db.Model(&model.UserProfile{}).
		Where("lower(username) = lower(?) AND user_id <> ?", profile.Username, profile.UserID).
		Count(&conflict).Error; err != nil {
		return err
	}
	if conflict == 0 && profile.Username != "" {
		existing.Username = profile.Username
	}
	if profile.Nickname != "" {
		existing.Nickname = profile.Nickname
	}
	if profile.AvatarURL != "" {
		existing.AvatarURL = profile.AvatarURL
	}
	return r.db.Save(&existing).Error
}

func (r *userProfileRepository) GetByUserID(userID string) (*model.UserProfile, error) {
	var profile model.UserProfile
	err := r.db.Where("user_id = ?", userID).First(&profile).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *userProfileRepository) GetByUsername(username string) (*model.UserProfile, error) {
	var profile model.UserProfile
	err := r.db.Where("lower(username) = lower(?)", username).First(&profile).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *userProfileRepository) ListByIDs(userIDs []string) (map[string]model.UserProfile, error) {
	result := make(map[string]model.UserProfile, len(userIDs))
	if len(userIDs) == 0 {
		return result, nil
	}
	var profiles []model.UserProfile
	if err := r.db.Where("user_id IN ?", userIDs).Find(&profiles).Error; err != nil {
		return nil, err
	}
	for _, profile := range profiles {
		result[profile.UserID] = profile
	}
	return result, nil
}
