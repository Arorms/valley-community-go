package repository

import (
	"errors"

	"gorm.io/gorm"

	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/model"
)

var ErrNotFound = errors.New("record not found")

type PostRepository interface {
	WithDB(db *gorm.DB) PostRepository
	Create(post *model.Post) error
	GetByID(id int64) (*model.Post, error)
	List(offset, limit int) ([]model.Post, error)
	ListByAuthor(authorID string, offset, limit int) ([]model.Post, error)
	Count() (int64, error)
	CountByAuthor(authorID string) (int64, error)
	IncrementCommentCount(postID int64, delta int) error
	IncrementLikeCount(postID int64, delta int) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) WithDB(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(post *model.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) GetByID(id int64) (*model.Post, error) {
	var post model.Post
	err := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&post).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) List(offset, limit int) ([]model.Post, error) {
	var posts []model.Post
	err := r.db.Where("is_deleted = ?", false).
		Order("created_at DESC, id DESC").
		Offset(offset).Limit(limit).
		Find(&posts).Error
	return posts, err
}

func (r *postRepository) ListByAuthor(authorID string, offset, limit int) ([]model.Post, error) {
	var posts []model.Post
	err := r.db.Where("author_id = ? AND is_deleted = ?", authorID, false).
		Order("created_at DESC, id DESC").
		Offset(offset).Limit(limit).
		Find(&posts).Error
	return posts, err
}

func (r *postRepository) Count() (int64, error) {
	var total int64
	err := r.db.Model(&model.Post{}).Where("is_deleted = ?", false).Count(&total).Error
	return total, err
}

func (r *postRepository) CountByAuthor(authorID string) (int64, error) {
	var total int64
	err := r.db.Model(&model.Post{}).
		Where("author_id = ? AND is_deleted = ?", authorID, false).
		Count(&total).Error
	return total, err
}

func (r *postRepository) IncrementCommentCount(postID int64, delta int) error {
	return r.db.Model(&model.Post{}).
		Where("id = ?", postID).
		UpdateColumn("comment_count", gorm.Expr("GREATEST(comment_count + ?, 0)", delta)).Error
}

func (r *postRepository) IncrementLikeCount(postID int64, delta int) error {
	return r.db.Model(&model.Post{}).
		Where("id = ?", postID).
		UpdateColumn("like_count", gorm.Expr("GREATEST(like_count + ?, 0)", delta)).Error
}
