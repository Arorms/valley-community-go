package viewmodel

import "gitlab.arorms.cn/arorms/applications/valley-community-go/internal/model"

type CurrentUser struct {
	Subject   string
	Username  string
	Nickname  string
	AvatarURL string
	Initials  string
	CSRFToken string
}

type UserSummary struct {
	UserID    string
	Username  string
	Nickname  string
	AvatarURL string
	Initials  string
}

func NewUserSummary(userID string, profile *model.UserProfile) UserSummary {
	summary := UserSummary{UserID: userID}
	if profile != nil {
		summary.Username = profile.Username
		summary.Nickname = profile.Nickname
		summary.AvatarURL = profile.AvatarURL
	}
	name := summary.Nickname
	if name == "" {
		name = summary.Username
	}
	summary.Initials = Initials(name)
	return summary
}

func Initials(name string) string {
	runes := []rune(name)
	if len(runes) == 0 {
		return "?"
	}
	if len(runes) == 1 {
		return string(runes[0])
	}
	return string(runes[:1])
}

func (s UserSummary) DisplayName() string {
	if s.Nickname != "" {
		return s.Nickname
	}
	if s.Username != "" {
		return s.Username
	}
	return "未知用户"
}
