package middlewares

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
)

type contextKey string

const userKey contextKey = "user"

// User 是校验通过后放入上下文的用户信息
type User struct {
	Subject           string `json:"subject"`
	Email             string `json:"email"`
	PreferredUsername string `json:"preferred_username"`
}

// tokenClaims 是 token payload 中本中间件关心的字段
type tokenClaims struct {
	Subject           string `json:"sub"`
	Email             string `json:"email"`
	PreferredUsername string `json:"preferred_username"`
	AuthorizedParty   string `json:"azp"`
}

// OIDCOptions 是构造 TokenValidator 所需的配置
type OIDCOptions struct {
	IssuerURL string
	Audience  string
	ClientID  string
}

type TokenValidator struct {
	verifier *oidc.IDTokenVerifier
	clientID string
}

// NewTokenValidator 初始化 Token 校验器
func NewTokenValidator(ctx context.Context, opts OIDCOptions) (*TokenValidator, error) {
	audience := opts.Audience
	if audience == "" {
		audience = opts.ClientID
	}
	if opts.IssuerURL == "" {
		return nil, errors.New("oidc: IssuerURL must be set")
	}
	if audience == "" {
		return nil, errors.New("oidc: Audience or ClientID must be set")
	}

	provider, err := oidc.NewProvider(ctx, opts.IssuerURL)
	if err != nil {
		return nil, err
	}

	return &TokenValidator{
		verifier: provider.Verifier(&oidc.Config{ClientID: audience}),
		clientID: opts.ClientID,
	}, nil
}

// Authenticate 中间件：仅校验 Token 是否存在且合法
func (v *TokenValidator) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. 检查 Header 是否存在
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			writeJSONError(w, http.StatusUnauthorized, "missing or invalid Authorization header")
			return
		}

		rawToken := strings.TrimPrefix(authHeader, "Bearer ")

		// 2. 校验 Token 合法性
		idToken, err := v.verifier.Verify(r.Context(), rawToken)
		if err != nil {
			writeJSONError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		// 3. 确认 token 是由预期的 client 申请的
		var rawClaims json.RawMessage
		if err := idToken.Claims(&rawClaims); err != nil {
			writeJSONError(w, http.StatusUnauthorized, "invalid token claims")
			return
		}
		claims, err := parseClaims(rawClaims)
		if err != nil {
			writeJSONError(w, http.StatusUnauthorized, "invalid token claims")
			return
		}
		if v.clientID != "" && claims.AuthorizedParty != "" && claims.AuthorizedParty != v.clientID {
			writeJSONError(w, http.StatusUnauthorized, "token was not issued for this client")
			return
		}

		// 4. 校验通过，将用户信息放入上下文传递给后续节点
		user := User{
			Subject:           claims.Subject,
			Email:             claims.Email,
			PreferredUsername: claims.PreferredUsername,
		}
		ctx := context.WithValue(r.Context(), userKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserFromContext 取出 Authenticate 写入上下文的用户信息
func UserFromContext(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(userKey).(User)
	return user, ok
}

// parseClaims 解析 token payload
func parseClaims(raw []byte) (tokenClaims, error) {
	var claims tokenClaims
	if err := json.Unmarshal(raw, &claims); err != nil {
		return tokenClaims{}, fmt.Errorf("parse token claims: %w", err)
	}
	return claims, nil
}

// writeJSONError 与 handlers.respondWithError 保持同样的响应结构
func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("WWW-Authenticate", `Bearer error="invalid_token"`)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
