package service

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/content"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/model"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/repository"
)

type CommentService interface {
	Create(postID int64, authorID string, parentID *int64, markdown string) (*model.Comment, error)
}

type commentService struct {
	db            *gorm.DB
	posts         repository.PostRepository
	comments      repository.CommentRepository
	profiles      repository.UserProfileRepository
	notifications repository.NotificationRepository
}

func NewCommentService(
	db *gorm.DB,
	posts repository.PostRepository,
	comments repository.CommentRepository,
	profiles repository.UserProfileRepository,
	notifications repository.NotificationRepository,
) CommentService {
	return &commentService{
		db:            db,
		posts:         posts,
		comments:      comments,
		profiles:      profiles,
		notifications: notifications,
	}
}

func (s *commentService) Create(postID int64, authorID string, parentID *int64, markdown string) (*model.Comment, error) {
	markdown = strings.TrimSpace(markdown)
	if markdown == "" {
		return nil, errors.New("comment is required")
	}
	if len([]rune(markdown)) > 10000 {
		return nil, errors.New("comment is too long")
	}
	html, err := content.Render(markdown)
	if err != nil {
		return nil, err
	}

	var created *model.Comment
	err = s.db.Transaction(func(tx *gorm.DB) error {
		post, err := s.posts.WithDB(tx).GetByID(postID)
		if err != nil {
			return err
		}

		var rootID *int64
		var replyToUserID string
		if parentID != nil {
			parent, err := s.comments.WithDB(tx).GetByID(*parentID)
			if err != nil {
				return err
			}
			if parent.PostID != postID {
				return errors.New("parent comment belongs to another post")
			}
			if parent.ParentID != nil {
				rootID = parent.ParentID
			} else {
				id := parent.ID
				rootID = &id
			}
			replyToUserID = parent.AuthorID
		}

		comment := &model.Comment{
			PostID:          postID,
			AuthorID:        authorID,
			ParentID:        rootID,
			ContentMarkdown: markdown,
			ContentHTML:     html,
		}
		if replyToUserID != "" {
			comment.ReplyToUserID = &replyToUserID
		}
		if err := s.comments.WithDB(tx).Create(comment); err != nil {
			return err
		}
		if err := s.posts.WithDB(tx).IncrementCommentCount(postID, 1); err != nil {
			return err
		}

		recipients := make(map[string]model.NotificationType)
		priorities := make(map[string]int)
		add := func(userID string, kind model.NotificationType, priority int) {
			if userID == "" || userID == authorID {
				return
			}
			if current, ok := priorities[userID]; ok && current >= priority {
				return
			}
			recipients[userID] = kind
			priorities[userID] = priority
		}

		if replyToUserID != "" {
			add(replyToUserID, model.NotificationTypeReply, 3)
		}
		add(post.AuthorID, model.NotificationTypeComment, 2)

		for _, username := range content.Mentions(markdown) {
			profile, err := s.profiles.GetByUsername(username)
			if errors.Is(err, repository.ErrNotFound) {
				continue
			}
			if err != nil {
				return err
			}
			add(profile.UserID, model.NotificationTypeMention, 1)
		}

		notifications := make([]model.Notification, 0, len(recipients))
		for userID, kind := range recipients {
			pid := postID
			cid := comment.ID
			notifications = append(notifications, model.Notification{
				RecipientID: userID,
				ActorID:     authorID,
				Type:        kind,
				PostID:      &pid,
				CommentID:   &cid,
			})
		}
		if err := s.notifications.WithDB(tx).CreateBatch(notifications); err != nil {
			return err
		}

		created = comment
		return nil
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}
