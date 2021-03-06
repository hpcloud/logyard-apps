package systail

import (
	"github.com/ActiveState/logyard-apps/common"
	"github.com/ActiveState/stackato-go/server"
)

type Config struct {
	MaxRecordSize int               `json:"max_record_size"`
	LogFiles      map[string]string `json:"log_files"`
}

var c *server.Config

func GetConfig() *Config {
	return c.GetConfig().(*Config)
}

func LoadConfig() {
	var err error
	c, err = server.NewConfig("systail", Config{})
	if err != nil {
		common.Fatal("Unable to load systail config; %v", err)
	}
}
