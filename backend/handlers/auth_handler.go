package handlers

import (
	"net/http"

	"valley-community-go/middlewares"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// Me
// GET /api/auth/me
// 返回当前 token 对应的用户信息，需要经过 Authenticate 中间件。
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	user, ok := middlewares.UserFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "unauthenticated")
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}