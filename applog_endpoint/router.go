package applog_endpoint

import (
	"github.com/ActiveState/log"
	"github.com/ActiveState/stackato-go/server"
	"github.com/apcera/nats"
	"strings"
)

var NATS *nats.EncodedConn

type routerRegisterInfo struct {
	Host string   `json:"host"`
	Port int      `json:"port"`
	URIs []string `json:"uris"`
	Tags struct {
		Component string `json:"string"`
	} `json:"tags"`
}

type applogEndpointConfigStruct struct {
	Hostname string `json:"hostname"`
}

func newRouterRegisterInfo() *routerRegisterInfo {
	var uri string
	applogEndpointConfig, err := server.NewConfig("applog_endpoint", applogEndpointConfigStruct{})
	if err == nil {
		uri = applogEndpointConfig.GetConfig().(*applogEndpointConfigStruct).Hostname
	}
	if uri == "" {
		clusterConfig := server.GetClusterConfig()
		uri = strings.Replace(clusterConfig.Endpoint, "api.", "logs.", 1)
	}

	info := new(routerRegisterInfo)
	info.Host = server.NodeIPMust()
	info.Port = PORT
	info.URIs = []string{uri}
	info.Tags.Component = COMPONENT
	return info
}

func routerAdvertise(m interface{}) {
	info := newRouterRegisterInfo()
	log.Infof("Advertising ourself to router: %+v (router.start? %+v)",
		info, m)
	NATS.Publish("router.register", info)
}

func RouterMain() {
	NATS = server.NewNatsClient(3)
	routerAdvertise(nil)
	NATS.Subscribe("router.start", routerAdvertise)
}
