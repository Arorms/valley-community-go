package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	sessionName = "vc_session"
	oauthName   = "vc_oauth"
	sessionAge  = 7 * 24 * 3600
	oauthAge    = 10 * 60
)

// User 是登录后写入 session 的最小身份信息。
type User struct {
	Subject  string
	Username string
	Email    string
	Name     string
	Picture  string
	IDToken  string
}

func (u User) DisplayName() string {
	if u.Name != "" {
		return u.Name
	}
	if u.Username != "" {
		return u.Username
	}
	return u.Subject
}

func (u User) Initials() string {
	name := u.DisplayName()
	runes := []rune(name)
	if len(runes) == 0 {
		return "?"
	}
	if len(runes) == 1 {
		return string(runes[0])
	}
	return string(runes[:1])
}

type SessionManager struct {
	store  *sessions.CookieStore
	secure bool
}

func NewSessionManager(secret string, secure bool) *SessionManager {
	if secret == "" {
		buf := make([]byte, 32)
		if _, err := rand.Read(buf); err != nil {
			panic(err)
		}
		secret = hex.EncodeToString(buf)
		log.Println("WARN: SESSION_SECRET 为空，已生成临时密钥；服务重启后登录态会失效")
	}
	// 统一派生 32 字节 key，避免用户提供任意长度 secret 时 securecookie 报错。
	sum := sha256.Sum256([]byte(secret))

	store := sessions.NewCookieStore(sum[:])
	store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   sessionAge,
	}
	return &SessionManager{store: store, secure: secure}
}

func (m *SessionManager) session(r *http.Request) (*sessions.Session, error) {
	return m.store.Get(r, sessionName)
}

func (m *SessionManager) oauthSession(r *http.Request) (*sessions.Session, error) {
	return m.store.Get(r, oauthName)
}

// SetUser 登录成功后写入 session，并确保存在 CSRF token。
func (m *SessionManager) SetUser(w http.ResponseWriter, r *http.Request, user User) error {
	s, err := m.session(r)
	if err != nil {
		return err
	}
	s.Values["sub"] = user.Subject
	s.Values["username"] = user.Username
	s.Values["email"] = user.Email
	s.Values["name"] = user.Name
	s.Values["picture"] = user.Picture
	s.Values["id_token"] = user.IDToken
	// 登录时轮换 CSRF token，避免会话固定。
	s.Values["csrf"] = RandomToken()
	return s.Save(r, w)
}

func (m *SessionManager) User(r *http.Request) (User, bool) {
	s, err := m.session(r)
	if err != nil {
		return User{}, false
	}
	sub, _ := s.Values["sub"].(string)
	if sub == "" {
		return User{}, false
	}
	user := User{Subject: sub}
	user.Username, _ = s.Values["username"].(string)
	user.Email, _ = s.Values["email"].(string)
	user.Name, _ = s.Values["name"].(string)
	user.Picture, _ = s.Values["picture"].(string)
	user.IDToken, _ = s.Values["id_token"].(string)
	return user, true
}

func (m *SessionManager) CSRFToken(r *http.Request) string {
	s, err := m.session(r)
	if err != nil {
		return ""
	}
	token, _ := s.Values["csrf"].(string)
	return token
}

func (m *SessionManager) ValidateCSRF(r *http.Request) bool {
	got := r.PostFormValue("_csrf")
	want := m.CSRFToken(r)
	if want == "" || got == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(got), []byte(want)) == 1
}

func (m *SessionManager) Clear(w http.ResponseWriter, r *http.Request) error {
	s, err := m.session(r)
	if err != nil {
		return err
	}
	s.Options.MaxAge = -1
	s.Values = map[interface{}]interface{}{}
	return s.Save(r, w)
}

func (m *SessionManager) SetOAuthState(w http.ResponseWriter, r *http.Request, state, nonce, verifier, returnTo string) error {
	s, err := m.oauthSession(r)
	if err != nil {
		return err
	}
	s.Options.MaxAge = oauthAge
	s.Values["state"] = state
	s.Values["nonce"] = nonce
	s.Values["verifier"] = verifier
	s.Values["return_to"] = returnTo
	return s.Save(r, w)
}

func (m *SessionManager) PopOAuthState(w http.ResponseWriter, r *http.Request) (state, nonce, verifier, returnTo string, err error) {
	s, err := m.oauthSession(r)
	if err != nil {
		return "", "", "", "", err
	}
	state, _ = s.Values["state"].(string)
	nonce, _ = s.Values["nonce"].(string)
	verifier, _ = s.Values["verifier"].(string)
	returnTo, _ = s.Values["return_to"].(string)

	s.Options.MaxAge = -1
	s.Values = map[interface{}]interface{}{}
	if saveErr := s.Save(r, w); saveErr != nil {
		return "", "", "", "", saveErr
	}
	return state, nonce, verifier, returnTo, nil
}

func (m *SessionManager) SetFlash(w http.ResponseWriter, r *http.Request, message string) {
	s, err := m.session(r)
	if err != nil {
		return
	}
	s.Values["flash"] = message
	_ = s.Save(r, w)
}

func (m *SessionManager) PopFlash(w http.ResponseWriter, r *http.Request) string {
	s, err := m.session(r)
	if err != nil {
		return ""
	}
	message, _ := s.Values["flash"].(string)
	if message == "" {
		return ""
	}
	delete(s.Values, "flash")
	_ = s.Save(r, w)
	return message
}

// RandomToken 生成 32 字节随机 token（hex）。
func RandomToken() string {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		panic(err)
	}
	return hex.EncodeToString(buf)
}
