package viewmodel

import "time"

type CommentView struct {
	ID          int64
	PostID      int64
	ParentID    *int64
	Author      UserSummary
	ReplyTo     *UserSummary
	ContentHTML string
	ContentText string
	CreatedAt   time.Time
	LikeCount   int
	Liked       bool
	Replies     []CommentView
}
