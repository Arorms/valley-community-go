package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/model"
)

type LikeRepository interface {
	WithDB(db *gorm.DB) LikeRepository
	LikePost(postID int64, userID string) (bool, error)
	UnlikePost(postID int64, userID string) (bool, error)
	LikeComment(commentID int64, userID string) (bool, error)
	UnlikeComment(commentID int64, userID string) (bool, error)
	LikedPosts(userID string, postIDs []int64) (map[int64]bool, error)
	LikedComments(userID string, commentIDs []int64) (map[int64]bool, error)
}

type likeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) LikeRepository {
	return &likeRepository{db: db}
}

func (r *likeRepository) WithDB(db *gorm.DB) LikeRepository {
	return &likeRepository{db: db}
}

func (r *likeRepository) LikePost(postID int64, userID string) (bool, error) {
	result := r.db.Clauses(clause.OnConflict{DoNothing: true}).
		Create(&model.PostLike{PostID: postID, UserID: userID})
	return result.RowsAffected > 0, result.Error
}

func (r *likeRepository) UnlikePost(postID int64, userID string) (bool, error) {
	result := r.db.Where("post_id = ? AND user_id = ?", postID, userID).
		Delete(&model.PostLike{})
	return result.RowsAffected > 0, result.Error
}

func (r *likeRepository) LikeComment(commentID int64, userID string) (bool, error) {
	result := r.db.Clauses(clause.OnConflict{DoNothing: true}).
		Create(&model.CommentLike{CommentID: commentID, UserID: userID})
	return result.RowsAffected > 0, result.Error
}

func (r *likeRepository) UnlikeComment(commentID int64, userID string) (bool, error) {
	result := r.db.Where("comment_id = ? AND user_id = ?", commentID, userID).
		Delete(&model.CommentLike{})
	return result.RowsAffected > 0, result.Error
}

func (r *likeRepository) LikedPosts(userID string, postIDs []int64) (map[int64]bool, error) {
	if len(postIDs) == 0 {
		return map[int64]bool{}, nil
	}
	var likes []model.PostLike
	err := r.db.Where("user_id = ? AND post_id IN ?", userID, postIDs).Find(&likes).Error
	if err != nil {
		return nil, err
	}
	result := make(map[int64]bool, len(likes))
	for _, like := range likes {
		result[like.PostID] = true
	}
	return result, nil
}

func (r *likeRepository) LikedComments(userID string, commentIDs []int64) (map[int64]bool, error) {
	if len(commentIDs) == 0 {
		return map[int64]bool{}, nil
	}
	var likes []model.CommentLike
	err := r.db.Where("user_id = ? AND comment_id IN ?", userID, commentIDs).Find(&likes).Error
	if err != nil {
		return nil, err
	}
	result := make(map[int64]bool, len(likes))
	for _, like := range likes {
		result[like.CommentID] = true
	}
	return result, nil
}
