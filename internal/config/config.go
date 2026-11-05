package config

import (
	"os"

	"github.com/jinzhu/configor"
)

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	DB      DBConfig      `yaml:"db"`
	OIDC    OIDCConfig    `yaml:"oidc"`
	Session SessionConfig `yaml:"session"`
}

type ServerConfig struct {
	Port    string `yaml:"port" default:"3000" env:"SERVER_PORT"`
	BaseURL string `yaml:"base_url" default:"http://localhost:3000" env:"SERVER_BASE_URL"`
}

type DBConfig struct {
	Host     string `yaml:"host" default:"localhost" env:"DB_HOST"`
	Port     string `yaml:"port" default:"5432" env:"DB_PORT"`
	User     string `yaml:"user" default:"postgres" env:"DB_USER"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	Database string `yaml:"database" default:"valley_db" env:"DB_NAME"`
}

type OIDCConfig struct {
	IssuerURL     string   `yaml:"issuer_url" env:"OIDC_ISSUER_URL"`
	ClientID      string   `yaml:"client_id" env:"OIDC_CLIENT_ID"`
	ClientSecret  string   `yaml:"client_secret" env:"OIDC_CLIENT_SECRET"`
	Audience      string   `yaml:"audience" env:"OIDC_AUDIENCE"`
	RedirectURL   string   `yaml:"redirect_url" env:"OIDC_REDIRECT_URL"`
	PostLogoutURL string   `yaml:"post_logout_redirect_url" env:"OIDC_POST_LOGOUT_REDIRECT_URL"`
	Scopes        []string `yaml:"scopes"`
}

type SessionConfig struct {
	Secret string `yaml:"secret" env:"SESSION_SECRET"`
}

func Load() (*Config, error) {
	path := "config.yml"
	if custom := os.Getenv("CONFIG_PATH"); custom != "" {
		path = custom
	}

	cfg := &Config{}
	if err := configor.Load(cfg, path); err != nil {
		return nil, err
	}
	return cfg, nil
}
