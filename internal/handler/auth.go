package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"

	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/auth"
)

type AuthHandler struct {
	*Base
	Client *auth.OIDCClient
}

func NewAuthHandler(base *Base, client *auth.OIDCClient) *AuthHandler {
	return &AuthHandler{Base: base, Client: client}
}

func (c *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	returnTo := safeReturnTo(r.URL.Query().Get("return_to"))
	state := auth.RandomToken()
	nonce := auth.RandomToken()
	verifier := oauth2.GenerateVerifier()

	if err := c.Sessions.SetOAuthState(w, r, state, nonce, verifier, returnTo); err != nil {
		c.renderError(w, r, http.StatusInternalServerError, "无法开始登录")
		return
	}
	http.Redirect(w, r, c.Client.AuthCodeURL(state, nonce, verifier), http.StatusFound)
}

func (c *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	if message := r.URL.Query().Get("error"); message != "" {
		c.renderError(w, r, http.StatusBadRequest, "登录失败："+message)
		return
	}

	state, nonce, verifier, returnTo, err := c.Sessions.PopOAuthState(w, r)
	if err != nil {
		c.renderError(w, r, http.StatusBadRequest, "登录状态已失效，请重试")
		return
	}
	if state == "" || r.URL.Query().Get("state") != state {
		c.renderError(w, r, http.StatusBadRequest, "登录状态校验失败")
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		c.renderError(w, r, http.StatusBadRequest, "缺少授权码")
		return
	}

	user, err := c.Client.Exchange(r.Context(), code, verifier, nonce)
	if err != nil {
		log.Printf("oidc exchange: %v", err)
		c.renderError(w, r, http.StatusUnauthorized, "登录校验失败")
		return
	}
	if err := c.Profiles.UpsertFromUser(user); err != nil {
		log.Printf("upsert user profile: %v", err)
	}
	if err := c.Sessions.SetUser(w, r, user); err != nil {
		c.renderError(w, r, http.StatusInternalServerError, "无法保存登录状态")
		return
	}

	http.Redirect(w, r, safeReturnTo(returnTo), http.StatusFound)
}

func (c *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	user, _ := c.Sessions.User(r)
	logoutURL := c.Client.LogoutURL(user.IDToken)
	if err := c.Sessions.Clear(w, r); err != nil {
		log.Printf("clear session: %v", err)
	}
	if logoutURL != "" {
		http.Redirect(w, r, logoutURL, http.StatusFound)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func safeReturnTo(value string) string {
	if strings.HasPrefix(value, "/") && !strings.HasPrefix(value, "//") {
		return value
	}
	return "/"
}

func pathParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}
