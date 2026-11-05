package auth

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/uuid"
	"golang.org/x/oauth2"

	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/config"
)

type OIDCClient struct {
	provider           *oidc.Provider
	oauth2Config       oauth2.Config
	verifier           *oidc.IDTokenVerifier
	endSessionEndpoint string
	postLogoutURL      string
}

type oidcClaims struct {
	Subject           string `json:"sub"`
	Email             string `json:"email"`
	PreferredUsername string `json:"preferred_username"`
	Name              string `json:"name"`
	Picture           string `json:"picture"`
}

func NewOIDCClient(ctx context.Context, cfg config.OIDCConfig, baseURL string) (*OIDCClient, error) {
	provider, err := oidc.NewProvider(ctx, cfg.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("oidc discovery: %w", err)
	}

	redirectURL := cfg.RedirectURL
	if redirectURL == "" {
		redirectURL = strings.TrimRight(baseURL, "/") + "/auth/callback"
	}
	scopes := cfg.Scopes
	if len(scopes) == 0 {
		scopes = []string{oidc.ScopeOpenID, "profile", "email"}
	}

	client := &OIDCClient{
		provider: provider,
		oauth2Config: oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			Endpoint:     provider.Endpoint(),
			RedirectURL:  redirectURL,
			Scopes:       scopes,
		},
		// ID token 的 aud 是 client_id；audience mapper 只作用于 access token。
		verifier:      provider.Verifier(&oidc.Config{ClientID: cfg.ClientID}),
		postLogoutURL: cfg.PostLogoutURL,
	}

	var discovery struct {
		EndSessionEndpoint string `json:"end_session_endpoint"`
	}
	if err := provider.Claims(&discovery); err == nil {
		client.endSessionEndpoint = discovery.EndSessionEndpoint
	}

	return client, nil
}

func (c *OIDCClient) AuthCodeURL(state, nonce, verifier string) string {
	return c.oauth2Config.AuthCodeURL(
		state,
		oidc.Nonce(nonce),
		oauth2.S256ChallengeOption(verifier),
	)
}

// Exchange 用授权码换取 token，并返回校验过的用户身份。
func (c *OIDCClient) Exchange(ctx context.Context, code, verifier, expectedNonce string) (User, error) {
	token, err := c.oauth2Config.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		return User{}, fmt.Errorf("exchange code: %w", err)
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok || rawIDToken == "" {
		return User{}, fmt.Errorf("token response has no id_token")
	}

	idToken, err := c.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return User{}, fmt.Errorf("verify id token: %w", err)
	}
	if idToken.Nonce != expectedNonce {
		return User{}, fmt.Errorf("id token nonce mismatch")
	}

	var claims oidcClaims
	if err := idToken.Claims(&claims); err != nil {
		return User{}, fmt.Errorf("parse id token claims: %w", err)
	}
	if _, err := uuid.Parse(claims.Subject); err != nil {
		return User{}, fmt.Errorf("oidc subject %q is not a uuid: %w", claims.Subject, err)
	}

	return User{
		Subject:  claims.Subject,
		Username: claims.PreferredUsername,
		Email:    claims.Email,
		Name:     claims.Name,
		Picture:  claims.Picture,
		IDToken:  rawIDToken,
	}, nil
}

// LogoutURL 返回 Keycloak 单点退出地址；未配置 post logout URI 时返回空。
func (c *OIDCClient) LogoutURL(idToken string) string {
	if c.endSessionEndpoint == "" || c.postLogoutURL == "" {
		return ""
	}
	values := url.Values{}
	values.Set("client_id", c.oauth2Config.ClientID)
	values.Set("post_logout_redirect_uri", c.postLogoutURL)
	if idToken != "" {
		values.Set("id_token_hint", idToken)
	}
	return c.endSessionEndpoint + "?" + values.Encode()
}
