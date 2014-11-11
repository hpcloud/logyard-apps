package util

import (
	"github.com/ActiveState/logyard-apps/common"
	"github.com/ActiveState/stackato-go/server"
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
	c, err = server.NewConfig("apptail", Config{})
	if err != nil {
		common.Fatal("Unable to load apptail config; %v", err)
	}
}

func GetBrandName() string {
	return getConfig().Info.Name
}

func init() {
	loadConfig()
}
