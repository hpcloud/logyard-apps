package config

import (
	"github.com/ActiveState/log"
)

type Config struct {
	DrainLifetime string `json:"drain_lifetime"`
	Hostname      string `json:"hostname"`
}

var c *ServerConfig

func GetConfig() *Config {
	return c.GetConfig().(*Config)
}

func LoadConfig() {
	var err error
	c, err = NewConfig("applog_endpoint", Config{})
	if err != nil {
		log.Fatal("Unable to load applog_endpoint config; %v", err)
	}
}
