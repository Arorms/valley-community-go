package handler

import (
	"net/http"
	"strconv"

	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/auth"
)

type CommentHandler struct {
	*Base
}

func NewCommentHandler(base *Base) *CommentHandler {
	return &CommentHandler{Base: base}
}

func (c *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	var parentID *int64
	if raw := r.FormValue("parent_id"); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			c.renderError(w, r, http.StatusBadRequest, "回复目标无效")
			return
		}
		parentID = &value
	}

	target := "/posts/" + strconv.FormatInt(postID, 10) + "#comments"
	if _, err := c.Comments.Create(postID, user.Subject, parentID, r.FormValue("content")); err != nil {
		c.redirectWithFlash(w, r, target, "评论失败："+err.Error())
		return
	}
	c.redirectWithFlash(w, r, target, "评论成功")
}
