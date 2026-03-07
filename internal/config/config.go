package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	APIPort    string `mapstructure:"API_PORT"`
	JWTSecret  string `mapstructure:"JWT_SECRET"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	UploadDir  string `mapstructure:"UPLOAD_DIR"`
}

func Load() (*Config, error) {
	viper.SetDefault("API_PORT", "8080")
	viper.SetDefault("JWT_SECRET", "change-me-in-production")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "3306")
	viper.SetDefault("DB_USER", "bobsbinder")
	viper.SetDefault("DB_PASSWORD", "bobsbinder")
	viper.SetDefault("DB_NAME", "bobsbinder")
	viper.SetDefault("UPLOAD_DIR", "./uploads")

	viper.AutomaticEnv()

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) DSN() string {
	return c.DBUser + ":" + c.DBPassword + "@tcp(" + c.DBHost + ":" + c.DBPort + ")/" + c.DBName + "?parseTime=true"
}
