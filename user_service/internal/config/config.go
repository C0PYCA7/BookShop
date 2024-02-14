package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string           `yaml:"env" envDefault:"local"`
	Database   DatabaseConfig   `yaml:"database" env-required:"true"`
	HttpServer HttpServerConfig `yaml:"http_server" env-required:"true"`
	Jwt        JwtConfig        `yaml:"jwt" env-required:"true"`
	Redis      RedisConfig      `yaml:"redis_database" env-required:"true"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type HttpServerConfig struct {
	Address     string        `yaml:"address" env-default:"localhost:8081"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type JwtConfig struct {
	SecretKey string `yaml:"secret_key"`
}

type RedisConfig struct {
	Addr string `yaml:"addr"`
}

func MustLoad() *Config {
	configPath := "user_service/config/config_user.yaml"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("config file doesn't exists: ", err)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("failed to read config: ", err)
	}

	return &cfg
}
