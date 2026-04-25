package config

import (
	"fmt"
	"web_socket/internal/utils"
)

type DataBaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLmode  string
}
type Config struct {
	DB            DataBaseConfig
	ServerAddress string
}

func NewConfig() *Config {
	addr := utils.GetEnv("SERVER_HOST", "localhost") + ":" + utils.GetEnv("SERVER_PORT", "8080")
	return &Config{
		DB: DataBaseConfig{
			Host:     utils.GetEnv("DB_HOST", "localhost"),
			Port:     utils.GetEnv("DB_PORT", "5432"),
			User:     utils.GetEnv("DB_USER", "postgres"),
			Password: utils.GetEnv("DB_PASSWORD", "postgres"),
			DBName:   utils.GetEnv("DB_NAME", "myapp"),
			SSLmode:  utils.GetEnv("DB_SSLMODE", "disable"),
		},
		ServerAddress: addr,
	}
}
func (c *Config) DNS() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", c.DB.Host, c.DB.Port, c.DB.User, c.DB.Password, c.DB.DBName, c.DB.SSLmode)
}
