package config

import (
	"time"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port         string `mapstructure:"port"`
	ReadTimeout  int32  `mapstructure:"read_timeout"`
	WriteTimeout int32  `mapstructure:"write_timeout"`
	MaxHeader    int32  `mapstructure:"max_header"`
}

type DatabaseConfig struct {
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Port     string `mapstructure:"port"`
}

type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`
	Expiration time.Duration `mapstructure:"expiration"`
}

type Config struct {
	SrvConfig  ServerConfig   `mapstructure:"server"`
	DBConfig   DatabaseConfig `mapstructure:"db"`
	JWTConfig  JWTConfig      `mapstructure:"jwt"`
	CtxTimeout time.Duration  `mapstructure:"timeout"`
}

func Read(configFilePath string) (*Config, error) {
	viper.AutomaticEnv()
	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
