package model

import "time"

type Comment struct {
	ID              int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	PostID          int64     `gorm:"column:post_id;not null" json:"post_id"`
	AuthorID        string    `gorm:"column:author_id;type:uuid;not null" json:"author_id"`
	ParentID        *int64    `gorm:"column:parent_id" json:"parent_id,omitempty"`
	ReplyToUserID   *string   `gorm:"column:reply_to_user_id;type:uuid" json:"reply_to_user_id,omitempty"`
	ContentMarkdown string    `gorm:"column:content_markdown;not null" json:"content_markdown"`
	ContentHTML     string    `gorm:"column:content_html;not null" json:"content_html"`
	LikeCount       int       `gorm:"column:like_count;not null;default:0" json:"like_count"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	IsDeleted       bool
}
