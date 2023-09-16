package util

import (
	"time"

	"github.com/spf13/viper"
)

// config stores all configurations of the application.
// the values are read by viper from a config file or environment variable
type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY=12345678901234567890123456789012"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION=15m"`
}

// LoadConfig reads configuration from the config file or environment variable
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
