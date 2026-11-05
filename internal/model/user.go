package model

import "time"

// UserProfile 是 Keycloak 用户资料的本地缓存，主键为 OIDC 的 sub。
// 登录/审计仍以 Keycloak 为准，这里只保存展示与 @ 解析所需的信息。
type UserProfile struct {
	UserID    string     `gorm:"column:user_id;type:uuid;primaryKey" json:"user_id"`
	Username  string     `gorm:"column:username;size:64;not null" json:"username"`
	Nickname  string     `gorm:"column:nickname;size:64" json:"nickname"`
	AvatarURL string     `gorm:"column:avatar_url" json:"avatar_url"`
	Gender    string     `gorm:"column:gender;size:16" json:"gender"`
	Birthday  *time.Time `gorm:"column:birthday;type:date" json:"birthday,omitempty"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}
