package config

import (
	"github.com/ActiveState/log"
	"stackato/server"
)

type Config struct {
	DrainLifetime string `json:"drain_lifetime"`
}

var c *server.Config

func GetConfig() *Config {
	return c.GetConfig().(*Config)
}

func LoadConfig() {
	var err error
	c, err = server.NewConfig("applog_endpoint", Config{})
	if err != nil {
		log.Fatal("Unable to load applog_endpoint config; %v", err)
	}
}
