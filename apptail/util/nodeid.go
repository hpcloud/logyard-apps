package util

import (
	"sync"

	"github.com/hpcloud/log"
	"github.com/hpcloud/logyard-apps/common"
	"github.com/hpcloud/stackato-go/server"
)

var once sync.Once
var nodeid string

// LocalNodeId returns the node ID of the local node.
func LocalNodeId() string {
	once.Do(func() {
		var err error
		nodeid, err = server.LocalIP()
		if err != nil {
			common.Fatal("Failed to determine IP addr: %v", err)
		}
		log.Info("Local Node ID: ", nodeid)
	})
	return nodeid
}
