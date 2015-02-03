// Mocked version of the shim to avoid attempting to make redis connections
// while running unit tests

// +build test

package config

import (
	"github.com/apcera/nats"
	"reflect"
	"sync"
)

type ConfDis struct {
	rootKey    string
	structType reflect.Type
	config     interface{} // Read-only view of current config tree.
	rev        int64
	mux        sync.Mutex  // Mutex to protect changes to config and rev.
	redis      interface{} // unused in mocks
	Changes    chan error  // Channel to receive config updates (value is return of reload())
}

type ServerConfig struct {
	name    string
	changes chan error
	*ConfDis
}

func (g *ServerConfig) SetConfig(config interface{}) {
	g.config = config
}

func (g *ServerConfig) GetConfig() interface{} {
	if g == nil {
		return nil
	}
	return g.config
}

type ClusterConfig struct {
	MbusIp   string `json:"mbusip"`
	Endpoint string `json:"endpoint"`
}

func NewConfig(group string, s interface{}) (*ServerConfig, error) {
	config := &ServerConfig{
		name:    group,
		changes: make(chan error),
		ConfDis: &ConfDis{
			rootKey: group,
			config:  s}}
	return config, nil
}

var clusterConfig *ClusterConfig
var clusterConfigInit sync.Once

func GetClusterConfig() *ClusterConfig {
	clusterConfigInit.Do(func() {
		clusterConfig = &ClusterConfig{
			MbusIp:   "127.0.0.1",
			Endpoint: "127.0.0.1"}
	})
	return clusterConfig
}

func NodeIPMust() string {
	return "127.0.0.1"
}

func NewNatsClient(retries int) *nats.EncodedConn {
	return nil
}

func init() {
	var err error
	c, err = NewConfig("a", &Config{})
	if err != nil {
		panic(err)
	}
}
