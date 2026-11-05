package viewmodel

import "time"

type PostCard struct {
	ID           int64
	Title        string
	Excerpt      string
	Author       UserSummary
	CreatedAt    time.Time
	LikeCount    int
	CommentCount int
	Liked        bool
	Favorited    bool
}

type PostDetail struct {
	PostCard
	HTMLContent string
}

type FeedPage struct {
	Page
	Posts      []PostCard
	PageNo     int
	TotalPages int
	Total      int64
}

type PostPage struct {
	Page
	Post         PostDetail
	Comments     []CommentView
	CommentCount int
}

type PostFormPage struct {
	Page
	FormTitle   string
	FormContent string
	Error       string
}
