package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

var EnvPrefix = ""
var RedisAddr = ":6379"

type (
	Config struct {
		*Server
	}

	Server struct {
		IP        string        `envconfig:"SERVER_IP" default:"127.0.0.1"`
		Port      string        `envconfig:"SERVER_PORT" default:"8088"`
		RedisPort string        `envconfig:"REDIS_PORT" default:":6379"`
		RTimeout  time.Duration `envconfig:"SERVER_READ_TIMEOUT" default:"15s"`
		WTimeout  time.Duration `envconfig:"SERVER_WRITE_TIMEOUT" default:"15s"`
	}
)

func LoadEnvConfig() *Config {
	var cfg Config
	if err := envconfig.Process(EnvPrefix, &cfg); err != nil {
		log.Fatalf("config: unable to load config for %T: %s", cfg, err)
	}
	return &cfg
}

func (s *Server) GetFullAddr() string {
	if s.Port == "" {
		return s.IP
	}
	return fmt.Sprintf("%s:%s", s.IP, s.Port)
}

func (s *Server) GetRedisPort() string {
	if s.RedisPort == "" {
		return s.RedisPort
	}
	return fmt.Sprintf("%s", s.RedisPort)
}
