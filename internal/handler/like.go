package handler

import (
	"net/http"
	"strconv"

	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/auth"
)

type LikeHandler struct {
	*Base
}

func NewLikeHandler(base *Base) *LikeHandler {
	return &LikeHandler{Base: base}
}

func (c *LikeHandler) TogglePost(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		c.renderError(w, r, http.StatusUnauthorized, "请先登录")
		return
	}
	postID, err := pathInt64(r, "id")
	if err != nil {
		c.renderError(w, r, http.StatusNotFound, "帖子不存在")
		return
	}
	fallback := "/posts/" + strconv.FormatInt(postID, 10)
	if _, err := c.Likes.TogglePost(postID, user.Subject); err != nil {
		c.redirectWithFlash(w, r, fallback, "操作失败")
		return
	}
	c.redirectBack(w, r, fallback)
}

func (c *LikeHandler) ToggleComment(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		c.renderError(w, r, http.StatusUnauthorized, "请先登录")
		return
	}
	commentID, err := pathInt64(r, "id")
	if err != nil {
		c.renderError(w, r, http.StatusNotFound, "评论不存在")
		return
	}
	fallback := "/"
	if postID := r.FormValue("post_id"); postID != "" {
		fallback = "/posts/" + postID + "#comment-" + strconv.FormatInt(commentID, 10)
	}
	if _, err := c.Likes.ToggleComment(commentID, user.Subject); err != nil {
		c.redirectWithFlash(w, r, fallback, "操作失败")
		return
	}
	c.redirectBack(w, r, fallback)
}
