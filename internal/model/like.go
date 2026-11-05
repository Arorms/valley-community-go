package model

import "time"

type PostLike struct {
	PostID    int64     `gorm:"column:post_id;primaryKey" json:"post_id"`
	UserID    string    `gorm:"column:user_id;type:uuid;primaryKey" json:"user_id"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

type CommentLike struct {
	CommentID int64     `gorm:"column:comment_id;primaryKey" json:"comment_id"`
	UserID    string    `gorm:"column:user_id;type:uuid;primaryKey" json:"user_id"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}
