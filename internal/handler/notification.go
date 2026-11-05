package handler

import (
	"net/http"

	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/auth"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/model/viewmodel"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/service"
)

type NotificationHandler struct {
	*Base
}

func NewNotificationHandler(base *Base) *NotificationHandler {
	return &NotificationHandler{Base: base}
}

func (c *NotificationHandler) Index(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		c.renderError(w, r, http.StatusUnauthorized, "请先登录")
		return
	}

	pageNo := queryPage(r)
	notifications, total, err := c.Notifications.List(user.Subject, pageNo, service.PageSize)
	if err != nil {
		c.renderError(w, r, http.StatusInternalServerError, "加载通知失败")
		return
	}

	c.render(w, http.StatusOK, "notifications.html", viewmodel.NotificationsPage{
		Page:          c.newPage(w, r, "通知", "notifications"),
		Notifications: notifications,
		PageNo:        pageNo,
		TotalPages:    service.TotalPages(total, service.PageSize),
		Total:         total,
	})
}

func (c *NotificationHandler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		c.renderError(w, r, http.StatusUnauthorized, "请先登录")
		return
	}
	if err := c.Notifications.MarkAllRead(user.Subject); err != nil {
		c.redirectWithFlash(w, r, "/notifications", "操作失败")
		return
	}
	c.redirectWithFlash(w, r, "/notifications", "已全部标记为已读")
}
