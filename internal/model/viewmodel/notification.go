package viewmodel

import (
	"time"

	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/model"
)

type NotificationView struct {
	ID        int64
	Type      model.NotificationType
	Actor     UserSummary
	Title     string
	Body      string
	Link      string
	Read      bool
	CreatedAt time.Time
}

type NotificationsPage struct {
	Page
	Notifications []NotificationView
	PageNo        int
	TotalPages    int
	Total         int64
}
