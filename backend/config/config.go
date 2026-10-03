package config

type Config struct {
	DB DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}
