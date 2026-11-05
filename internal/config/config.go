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
