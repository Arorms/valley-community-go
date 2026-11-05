package handler

import (
	"net/http"

	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/model/viewmodel"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/service"
)

type HomeHandler struct {
	*Base
}

func NewHomeHandler(base *Base) *HomeHandler {
	return &HomeHandler{Base: base}
}

func (c *HomeHandler) Index(w http.ResponseWriter, r *http.Request) {
	pageNo := queryPage(r)
	posts, total, err := c.Posts.List(pageNo, service.PageSize, c.viewerID(r))
	if err != nil {
		c.renderError(w, r, http.StatusInternalServerError, "加载帖子失败")
		return
	}

	c.render(w, http.StatusOK, "feed.html", viewmodel.FeedPage{
		Page:       c.newPage(w, r, "首页", "home"),
		Posts:      posts,
		PageNo:     pageNo,
		TotalPages: service.TotalPages(total, service.PageSize),
		Total:      total,
	})
}
