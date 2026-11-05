package model

import "time"

type Post struct {
	ID              int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	AuthorID        string    `gorm:"column:author_id;type:uuid;not null" json:"author_id"`
	Title           string    `gorm:"column:title;size:200;not null" json:"title"`
	MarkdownContent string    `gorm:"column:markdown_content;not null" json:"markdown_content"`
	HTMLContent     string    `gorm:"column:html_content;not null" json:"html_content"`
	LikeCount       int       `gorm:"column:like_count;not null;default:0" json:"like_count"`
	CommentCount    int       `gorm:"column:comment_count;not null;default:0" json:"comment_count"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	IsDeleted       bool
}
