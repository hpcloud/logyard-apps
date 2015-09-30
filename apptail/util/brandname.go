package util

import (
	"github.com/hpcloud/logyard-apps/common"
	"github.com/hpcloud/stackato-go/server"
)

type Config struct {
	Info struct {
		Name string `json:"name"`
	} `json:"info"`
}

var c *server.Config

func getConfig() *Config {
	return c.GetConfig().(*Config)
}

func loadConfig() {
	var err error
	c, err = server.NewConfig("cloud_controller_ng", Config{})
	if err != nil {
		common.Fatal("Unable to load cc_ng config; %v", err)
	}
}

func GetBrandName() string {
	return getConfig().Info.Name
}
