package model

type Favorite struct {
	UserID string `gorm:"column:user_id;type:uuid;primaryKey" json:"user_id"`
	PostID int64  `gorm:"column:post_id;primaryKey" json:"post_id"`
}
