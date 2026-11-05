package handler

import (
	"net/http"
	"strconv"

	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/auth"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/model/request"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/model/viewmodel"
)

type PostHandler struct {
	*Base
}

func NewPostHandler(base *Base) *PostHandler {
	return &PostHandler{Base: base}
}

func (c *PostHandler) New(w http.ResponseWriter, r *http.Request) {
	c.render(w, http.StatusOK, "post_new.html", viewmodel.PostFormPage{
		Page: c.newPage(w, r, "发帖", "home"),
	})
}

func (c *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		c.renderError(w, r, http.StatusUnauthorized, "请先登录")
		return
	}

	req := request.PostCreateRequest{
		Title:           r.FormValue("title"),
		MarkdownContent: r.FormValue("content"),
	}
	post, err := c.Posts.Create(user.Subject, req)
	if err != nil {
		c.render(w, http.StatusBadRequest, "post_new.html", viewmodel.PostFormPage{
			Page:        c.newPage(w, r, "发帖", "home"),
			FormTitle:   req.Title,
			FormContent: req.MarkdownContent,
			Error:       err.Error(),
		})
		return
	}

	c.redirectWithFlash(w, r, "/posts/"+strconv.FormatInt(post.ID, 10), "发布成功")
}

func (c *PostHandler) Show(w http.ResponseWriter, r *http.Request) {
	id, err := pathInt64(r, "id")
	if err != nil {
		c.renderError(w, r, http.StatusNotFound, "帖子不存在")
		return
	}

	detail, comments, err := c.Posts.Detail(id, c.viewerID(r))
	if isNotFound(err) {
		c.renderError(w, r, http.StatusNotFound, "帖子不存在或已删除")
		return
	}
	if err != nil {
		c.renderError(w, r, http.StatusInternalServerError, "加载帖子失败")
		return
	}

	c.render(w, http.StatusOK, "post.html", viewmodel.PostPage{
		Page:         c.newPage(w, r, detail.Title, "home"),
		Post:         *detail,
		Comments:     comments,
		CommentCount: detail.CommentCount,
	})
}
