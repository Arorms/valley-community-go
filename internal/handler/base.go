package handler

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/auth"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/model/viewmodel"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/repository"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/service"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/view"
)

type Base struct {
	Renderer      *view.Renderer
	Sessions      *auth.SessionManager
	Posts         service.PostService
	Comments      service.CommentService
	Likes         service.LikeService
	Favorites     service.FavoriteService
	Notifications service.NotificationService
	Profiles      service.UserProfileService
}

func (b *Base) newPage(w http.ResponseWriter, r *http.Request, title, active string) viewmodel.Page {
	page := viewmodel.Page{
		Title:  title,
		Active: active,
		Flash:  b.Sessions.PopFlash(w, r),
	}
	if user, ok := auth.UserFromContext(r.Context()); ok {
		page.User = b.currentUser(r, user)
		if count, err := b.Notifications.UnreadCount(user.Subject); err == nil {
			page.Unread = int(count)
		}
	}
	return page
}

func (b *Base) currentUser(r *http.Request, user auth.User) *viewmodel.CurrentUser {
	current := &viewmodel.CurrentUser{
		Subject:   user.Subject,
		Username:  user.Username,
		Nickname:  user.Name,
		AvatarURL: user.Picture,
		CSRFToken: b.Sessions.CSRFToken(r),
	}
	if profile, err := b.Profiles.GetByUserID(user.Subject); err == nil {
		current.Username = profile.Username
		current.Nickname = profile.Nickname
		current.AvatarURL = profile.AvatarURL
	}
	name := current.Nickname
	if name == "" {
		name = current.Username
	}
	current.Initials = viewmodel.Initials(name)
	return current
}

func (b *Base) viewerID(r *http.Request) string {
	if user, ok := auth.UserFromContext(r.Context()); ok {
		return user.Subject
	}
	return ""
}

func (b *Base) render(w http.ResponseWriter, status int, page string, data any) {
	var buf bytes.Buffer
	if err := b.Renderer.Render(&buf, page, data); err != nil {
		log.Printf("render %s: %v", page, err)
		http.Error(w, "render error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(buf.Bytes())
}

func (b *Base) renderError(w http.ResponseWriter, r *http.Request, status int, message string) {
	b.render(w, status, "error.html", viewmodel.ErrorPage{
		Page:    b.newPage(w, r, "出错了", ""),
		Message: message,
	})
}

func (b *Base) redirectBack(w http.ResponseWriter, r *http.Request, fallback string) {
	target := r.FormValue("return_to")
	if !strings.HasPrefix(target, "/") || strings.HasPrefix(target, "//") {
		target = fallback
	}
	http.Redirect(w, r, target, http.StatusSeeOther)
}

func (b *Base) redirectWithFlash(w http.ResponseWriter, r *http.Request, target, message string) {
	b.Sessions.SetFlash(w, r, message)
	http.Redirect(w, r, target, http.StatusSeeOther)
}

func queryPage(r *http.Request) int {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		return 1
	}
	return page
}

func pathInt64(r *http.Request, key string) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, key), 10, 64)
}

func isNotFound(err error) bool {
	return errors.Is(err, repository.ErrNotFound)
}
