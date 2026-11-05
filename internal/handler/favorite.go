package handler

import (
	"net/http"
	"strconv"

	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/auth"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/model/viewmodel"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/service"
)

type FavoriteHandler struct {
	*Base
}

func NewFavoriteHandler(base *Base) *FavoriteHandler {
	return &FavoriteHandler{Base: base}
}

func (c *FavoriteHandler) Toggle(w http.ResponseWriter, r *http.Request) {
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
	favorited, err := c.Favorites.Toggle(user.Subject, postID)
	if err != nil {
		c.redirectWithFlash(w, r, fallback, "操作失败")
		return
	}
	message := "已取消收藏"
	if favorited {
		message = "已加入收藏"
	}
	c.redirectWithFlash(w, r, fallback, message)
}

func (c *FavoriteHandler) List(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		c.renderError(w, r, http.StatusUnauthorized, "请先登录")
		return
	}
	pageNo := queryPage(r)
	posts, total, err := c.Favorites.List(user.Subject, pageNo, service.PageSize)
	if err != nil {
		c.renderError(w, r, http.StatusInternalServerError, "加载收藏失败")
		return
	}
	c.render(w, http.StatusOK, "favorites.html", viewmodel.FavoritesPage{
		Page:       c.newPage(w, r, "我的收藏", "favorites"),
		Posts:      posts,
		PageNo:     pageNo,
		TotalPages: service.TotalPages(total, service.PageSize),
		Total:      total,
	})
}
