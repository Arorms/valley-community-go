package handler

import (
	"net/http"

	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/model/viewmodel"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/service"
)

type ProfileHandler struct {
	*Base
}

func NewProfileHandler(base *Base) *ProfileHandler {
	return &ProfileHandler{Base: base}
}

func (c *ProfileHandler) Show(w http.ResponseWriter, r *http.Request) {
	username := pathParam(r, "username")
	profile, err := c.Profiles.GetByUsername(username)
	if isNotFound(err) {
		c.renderError(w, r, http.StatusNotFound, "用户不存在")
		return
	}
	if err != nil {
		c.renderError(w, r, http.StatusInternalServerError, "加载用户失败")
		return
	}

	pageNo := queryPage(r)
	posts, total, err := c.Posts.ListByAuthor(profile.UserID, pageNo, service.PageSize, c.viewerID(r))
	if err != nil {
		c.renderError(w, r, http.StatusInternalServerError, "加载帖子失败")
		return
	}

	summary := viewmodel.NewUserSummary(profile.UserID, profile)
	c.render(w, http.StatusOK, "profile.html", viewmodel.ProfilePage{
		Page:       c.newPage(w, r, summary.DisplayName(), "profile"),
		Profile:    summary,
		Posts:      posts,
		PageNo:     pageNo,
		TotalPages: service.TotalPages(total, service.PageSize),
		Total:      total,
		IsMe:       profile.UserID == c.viewerID(r),
	})
}
