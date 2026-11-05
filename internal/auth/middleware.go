package auth

import (
	"context"
	"net/http"
	"net/url"
)

type userContextKey struct{}

func WithUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, userContextKey{}, user)
}

func UserFromContext(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(userContextKey{}).(User)
	return user, ok
}

// LoadUser 每个请求从 session 读取用户并放进 context。
func (m *SessionManager) LoadUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user, ok := m.User(r); ok {
			r = r.WithContext(WithUser(r.Context(), user))
		}
		next.ServeHTTP(w, r)
	})
}

// RequireAuth 未登录时跳转登录页，并带回跳地址。
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := UserFromContext(r.Context()); !ok {
			returnTo := r.URL.RequestURI()
			http.Redirect(w, r, "/auth/login?return_to="+url.QueryEscape(returnTo), http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// CSRF 校验所有写操作的表单 token。
func (m *SessionManager) CSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
			if !m.ValidateCSRF(r) {
				http.Error(w, "invalid csrf token", http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
