package service

import (
	"errors"
	"strconv"

	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/content"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/model/viewmodel"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/repository"
)

type NotificationService interface {
	List(recipientID string, pageNo, pageSize int) ([]viewmodel.NotificationView, int64, error)
	UnreadCount(recipientID string) (int64, error)
	MarkAllRead(recipientID string) error
}

type notificationService struct {
	notifications repository.NotificationRepository
	posts         repository.PostRepository
	comments      repository.CommentRepository
	profiles      repository.UserProfileRepository
}

func NewNotificationService(
	notifications repository.NotificationRepository,
	posts repository.PostRepository,
	comments repository.CommentRepository,
	profiles repository.UserProfileRepository,
) NotificationService {
	return &notificationService{
		notifications: notifications,
		posts:         posts,
		comments:      comments,
		profiles:      profiles,
	}
}

func (s *notificationService) UnreadCount(recipientID string) (int64, error) {
	return s.notifications.CountUnread(recipientID)
}

func (s *notificationService) MarkAllRead(recipientID string) error {
	return s.notifications.MarkAllRead(recipientID)
}

func (s *notificationService) List(recipientID string, pageNo, pageSize int) ([]viewmodel.NotificationView, int64, error) {
	pageNo, pageSize = normalizePage(pageNo, pageSize)
	offset := (pageNo - 1) * pageSize

	notifications, err := s.notifications.ListByRecipient(recipientID, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.notifications.CountByRecipient(recipientID)
	if err != nil {
		return nil, 0, err
	}
	if len(notifications) == 0 {
		return []viewmodel.NotificationView{}, total, nil
	}

	actorIDs := make([]string, 0, len(notifications))
	for _, notification := range notifications {
		actorIDs = append(actorIDs, notification.ActorID)
	}
	profiles, err := s.profiles.ListByIDs(actorIDs)
	if err != nil {
		return nil, 0, err
	}

	views := make([]viewmodel.NotificationView, 0, len(notifications))
	for _, notification := range notifications {
		view := viewmodel.NotificationView{
			ID:        notification.ID,
			Type:      notification.Type,
			Actor:     viewmodel.NewUserSummary(notification.ActorID, profilePointer(profiles, notification.ActorID)),
			Read:      notification.ReadAt != nil,
			CreatedAt: notification.CreatedAt,
		}

		if notification.PostID != nil {
			post, err := s.posts.GetByID(*notification.PostID)
			if err == nil {
				view.Title = post.Title
				view.Link = "/posts/" + strconv.FormatInt(post.ID, 10)
			} else if errors.Is(err, repository.ErrNotFound) {
				view.Title = "帖子已删除"
				view.Link = "/"
			} else {
				return nil, 0, err
			}
		}

		if notification.CommentID != nil {
			comment, err := s.comments.GetByID(*notification.CommentID)
			if err == nil {
				view.Body = content.Excerpt(comment.ContentMarkdown, 120)
				if view.Link != "" {
					view.Link += "#comment-" + strconv.FormatInt(comment.ID, 10)
				}
			} else if !errors.Is(err, repository.ErrNotFound) {
				return nil, 0, err
			}
		}

		views = append(views, view)
	}
	return views, total, nil
}
