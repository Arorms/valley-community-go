package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/model"
)

type FavoriteRepository interface {
	WithDB(db *gorm.DB) FavoriteRepository
	Add(userID string, postID int64) (bool, error)
	Remove(userID string, postID int64) (bool, error)
	Exists(userID string, postID int64) (bool, error)
	FavoritedPosts(userID string, postIDs []int64) (map[int64]bool, error)
	ListPosts(userID string, offset, limit int) ([]model.Post, error)
	Count(userID string) (int64, error)
}

type favoriteRepository struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) FavoriteRepository {
	return &favoriteRepository{db: db}
}

func (r *favoriteRepository) WithDB(db *gorm.DB) FavoriteRepository {
	return &favoriteRepository{db: db}
}

func (r *favoriteRepository) Add(userID string, postID int64) (bool, error) {
	result := r.db.Clauses(clause.OnConflict{DoNothing: true}).
		Create(&model.Favorite{UserID: userID, PostID: postID})
	return result.RowsAffected > 0, result.Error
}

func (r *favoriteRepository) Remove(userID string, postID int64) (bool, error) {
	result := r.db.Where("user_id = ? AND post_id = ?", userID, postID).
		Delete(&model.Favorite{})
	return result.RowsAffected > 0, result.Error
}

func (r *favoriteRepository) Exists(userID string, postID int64) (bool, error) {
	var count int64
	err := r.db.Model(&model.Favorite{}).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Count(&count).Error
	return count > 0, err
}

func (r *favoriteRepository) FavoritedPosts(userID string, postIDs []int64) (map[int64]bool, error) {
	if len(postIDs) == 0 {
		return map[int64]bool{}, nil
	}
	var favorites []model.Favorite
	err := r.db.Where("user_id = ? AND post_id IN ?", userID, postIDs).Find(&favorites).Error
	if err != nil {
		return nil, err
	}
	result := make(map[int64]bool, len(favorites))
	for _, favorite := range favorites {
		result[favorite.PostID] = true
	}
	return result, nil
}

func (r *favoriteRepository) ListPosts(userID string, offset, limit int) ([]model.Post, error) {
	var posts []model.Post
	err := r.db.Model(&model.Post{}).
		Joins("JOIN favorites ON favorites.post_id = posts.id").
		Where("favorites.user_id = ? AND posts.is_deleted = ?", userID, false).
		Order("posts.created_at DESC, posts.id DESC").
		Offset(offset).Limit(limit).
		Find(&posts).Error
	return posts, err
}

func (r *favoriteRepository) Count(userID string) (int64, error) {
	var total int64
	err := r.db.Model(&model.Favorite{}).Where("user_id = ?", userID).Count(&total).Error
	return total, err
}
