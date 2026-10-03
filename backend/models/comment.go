package models

type Comment struct {
	//gorm.Model
	PostID  uint   `json:"post_id"`
	UserID  string `json:"user_id"`
	Content string `json:"content"`
}
