package repository

import (
	"errors"

	"gorm.io/gorm"

	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/model"
)

type CommentRepository interface {
	WithDB(db *gorm.DB) CommentRepository
	Create(comment *model.Comment) error
	GetByID(id int64) (*model.Comment, error)
	ListByPost(postID int64) ([]model.Comment, error)
	CountByPost(postID int64) (int64, error)
	IncrementLikeCount(commentID int64, delta int) error
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) WithDB(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(comment *model.Comment) error {
	return r.db.Create(comment).Error
}

func (r *commentRepository) GetByID(id int64) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&comment).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) ListByPost(postID int64) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.db.Where("post_id = ? AND is_deleted = ?", postID, false).
		Order("created_at ASC, id ASC").
		Find(&comments).Error
	return comments, err
}

func (r *commentRepository) CountByPost(postID int64) (int64, error) {
	var total int64
	err := r.db.Model(&model.Comment{}).
		Where("post_id = ? AND is_deleted = ?", postID, false).
		Count(&total).Error
	return total, err
}

func (r *commentRepository) IncrementLikeCount(commentID int64, delta int) error {
	return r.db.Model(&model.Comment{}).
		Where("id = ?", commentID).
		UpdateColumn("like_count", gorm.Expr("GREATEST(like_count + ?, 0)", delta)).Error
}
