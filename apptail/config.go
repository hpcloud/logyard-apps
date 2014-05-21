package apptail

import (
	"github.com/ActiveState/log"
	"github.com/ActiveState/logyard-apps/common"
	"github.com/ActiveState/tail/ratelimiter"
	"stackato/server"
	"time"
)

type leakyBucketConfig struct {
}

type Config struct {
	MaxRecordSize     int    `json:"max_record_size"`
	MaxLinesPerSecond int64  `json:"max_line_per_second"`
	MaxLinesBurst     uint16 `json:"max_lines_burst"`
	FileSizeLimit     int64  `json:"read_limit"`
}

var c *server.Config

func GetConfig() *Config {
	return c.GetConfig().(*Config)
}

func (c *Config) GetLeakyBucket() *ratelimiter.LeakyBucket {
	rate := c.MaxLinesPerSecond
	if rate < 1 {
		log.Warnf("max_lines_per_second must be a positive integer; using default")
		rate = 100
	}

	burstSize := c.MaxLinesBurst
	if burstSize < 1 {
		log.Warnf("max_lines_burst must be a positive integer; using default")
		burstSize = 10000
	}

	interval := time.Duration(int64(time.Second) / rate)

	return ratelimiter.NewLeakyBucket(burstSize, interval)
}

func LoadConfig() {
	var err error
	c, err = server.NewConfig("apptail", Config{})
	if err != nil {
		common.Fatal("Unable to load apptail config; %v", err)
	}
}
