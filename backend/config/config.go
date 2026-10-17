package config

import (
	"os"

	"github.com/jinzhu/configor"
)

type Config struct {
	DB   DBConfig   `yaml:"db"`
	OIDC OIDCConfig `yaml:"oidc"`
}

type DBConfig struct {
	Host     string `yaml:"host" default:"localhost" env:"DB_HOST"`
	Port     string `yaml:"port" default:"5432" env:"DB_PORT"`
	User     string `yaml:"user" default:"postgres" env:"DB_USER"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	Database string `yaml:"database" default:"valley_db" env:"DB_NAME"`
}

type OIDCConfig struct {
	IssuerURL string `yaml:"issuer_url"`
	Audience  string `yaml:"audience"`
	ClientID  string `yaml:"client_id"`
}

// Load 读取 config.yml，路径可用 CONFIG_PATH 环境变量覆盖。
//
// 注意 configor 的行为：指定的文件不存在时会依次尝试 config.<env>.yml、
// config.example.yml，都没有则只打印一行警告，然后用 default tag + 环境变量兜底。
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
