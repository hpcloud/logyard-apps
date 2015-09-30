// Shim for github.com/hpcloud/stackato-go/server to make it easier to
// replace during testing.

// +build !test

package config

import (
	"github.com/apcera/nats"
	"github.com/hpcloud/stackato-go/server"
)

type ServerConfig struct {
	*server.Config
}

func NewConfig(group string, s interface{}) (*ServerConfig, error) {
	config, err := server.NewConfig(group, s)
	return &ServerConfig{config}, err
}

type ClusterConfig struct {
	*server.ClusterConfig
}

func GetClusterConfig() *ClusterConfig {
	return &ClusterConfig{server.GetClusterConfig()}
}

func NodeIPMust() string {
	return server.NodeIPMust()
}

func NewNatsClient(retries int) *nats.EncodedConn {
	return server.NewNatsClient(retries)
}
