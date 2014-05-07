package main

import (
	"github.com/ActiveState/log"
	"stackato/server"
	"strings"
)

const COMPONENT = "websocket_endpoint"
const PORT = 5722

type routerRegisterInfo struct {
	Host string   `json:"host"`
	Port int      `json:"port"`
	URIs []string `json:"uris"`
	Tags struct {
		Component string `json:"string"`
	} `json:"tags"`
}

func newRouterRegisterInfo() *routerRegisterInfo {
	clusterConfig := server.GetClusterConfig()
	uri := strings.Replace(clusterConfig.Endpoint, "api.", "logs.", 1)

	info := new(routerRegisterInfo)
	info.Host = server.NodeIPMust()
	info.Port = PORT
	info.URIs = []string{uri}
	info.Tags.Component = COMPONENT
	return info
}

func advertiseToStackatoRouter() error {
	// Stackato work
	nats := server.NewNatsClient(3)

	info := newRouterRegisterInfo()
	log.Infof("Advertising ourself to router: %+v", info)
	return nats.Publish("router.register", info)
}
