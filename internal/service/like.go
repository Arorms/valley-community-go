package service

import (
	"gorm.io/gorm"

	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/model"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/repository"
)

type LikeService interface {
	TogglePost(postID int64, userID string) (bool, error)
	ToggleComment(commentID int64, userID string) (bool, error)
}

type likeService struct {
	db            *gorm.DB
	posts         repository.PostRepository
	comments      repository.CommentRepository
	likes         repository.LikeRepository
	notifications repository.NotificationRepository
}

func NewLikeService(
	db *gorm.DB,
	posts repository.PostRepository,
	comments repository.CommentRepository,
	likes repository.LikeRepository,
	notifications repository.NotificationRepository,
) LikeService {
	return &likeService{
		db:            db,
		posts:         posts,
		comments:      comments,
		likes:         likes,
		notifications: notifications,
	}
}

func (s *likeService) TogglePost(postID int64, userID string) (bool, error) {
	liked := false
	err := s.db.Transaction(func(tx *gorm.DB) error {
		post, err := s.posts.WithDB(tx).GetByID(postID)
		if err != nil {
			return err
		}

		created, err := s.likes.WithDB(tx).LikePost(postID, userID)
		if err != nil {
			return err
		}
		if created {
			liked = true
			if err := s.posts.WithDB(tx).IncrementLikeCount(postID, 1); err != nil {
				return err
			}
			if post.AuthorID != userID {
				pid := postID
				notification := model.Notification{
					RecipientID: post.AuthorID,
					ActorID:     userID,
					Type:        model.NotificationTypePostLike,
					PostID:      &pid,
				}
				if err := s.notifications.WithDB(tx).CreateBatch([]model.Notification{notification}); err != nil {
					return err
				}
			}
			return nil
		}

		removed, err := s.likes.WithDB(tx).UnlikePost(postID, userID)
		if err != nil {
			return err
		}
		if removed {
			return s.posts.WithDB(tx).IncrementLikeCount(postID, -1)
		}
		return nil
	})
	return liked, err
}

func (s *likeService) ToggleComment(commentID int64, userID string) (bool, error) {
	liked := false
	err := s.db.Transaction(func(tx *gorm.DB) error {
		comment, err := s.comments.WithDB(tx).GetByID(commentID)
		if err != nil {
			return err
		}

		created, err := s.likes.WithDB(tx).LikeComment(commentID, userID)
		if err != nil {
			return err
		}
		if created {
			liked = true
			if err := s.comments.WithDB(tx).IncrementLikeCount(commentID, 1); err != nil {
				return err
			}
			if comment.AuthorID != userID {
				pid := comment.PostID
				cid := commentID
				notification := model.Notification{
					RecipientID: comment.AuthorID,
					ActorID:     userID,
					Type:        model.NotificationTypeCommentLike,
					PostID:      &pid,
					CommentID:   &cid,
				}
				if err := s.notifications.WithDB(tx).CreateBatch([]model.Notification{notification}); err != nil {
					return err
				}
			}
			return nil
		}

		removed, err := s.likes.WithDB(tx).UnlikeComment(commentID, userID)
		if err != nil {
			return err
		}
		if removed {
			return s.comments.WithDB(tx).IncrementLikeCount(commentID, -1)
		}
		return nil
	})
	return liked, err
}
