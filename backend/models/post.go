package models

type Post struct {
	//gorm.Model
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  string `json:"user_id"`
	//Comments []Comment `json:"comments" gorm:"constraint:OnDelete:CASCADE;"`
}

type PostCreateRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
