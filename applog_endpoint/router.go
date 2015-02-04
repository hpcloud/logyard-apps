package applog_endpoint

import (
	"github.com/ActiveState/log"
	"github.com/ActiveState/logyard-apps/applog_endpoint/config"
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

func getApplogEndpointUri() string {
	uri := config.GetConfig().Hostname
	if uri == "" {
		clusterConfig := config.GetClusterConfig()
		uri = strings.Replace(clusterConfig.Endpoint, "api.", "logs.", 1)
	}
	return uri
}

func newRouterRegisterInfo() *routerRegisterInfo {
	info := new(routerRegisterInfo)
	info.Host = config.NodeIPMust()
	info.Port = PORT
	info.URIs = []string{getApplogEndpointUri()}
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
	NATS = config.NewNatsClient(3)
	routerAdvertise(nil)
	NATS.Subscribe("router.start", routerAdvertise)
}
