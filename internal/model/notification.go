package model

import "time"

type NotificationType string

const (
	NotificationTypeComment     NotificationType = "comment"
	NotificationTypeReply       NotificationType = "reply"
	NotificationTypeMention     NotificationType = "mention"
	NotificationTypePostLike    NotificationType = "post_like"
	NotificationTypeCommentLike NotificationType = "comment_like"
)

type Notification struct {
	ID          int64            `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RecipientID string           `gorm:"column:recipient_id;type:uuid;not null" json:"recipient_id"`
	ActorID     string           `gorm:"column:actor_id;type:uuid;not null" json:"actor_id"`
	Type        NotificationType `gorm:"column:type;size:32;not null" json:"type"`
	PostID      *int64           `gorm:"column:post_id" json:"post_id,omitempty"`
	CommentID   *int64           `gorm:"column:comment_id" json:"comment_id,omitempty"`
	ReadAt      *time.Time       `gorm:"column:read_at" json:"read_at,omitempty"`
	CreatedAt   time.Time        `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}
