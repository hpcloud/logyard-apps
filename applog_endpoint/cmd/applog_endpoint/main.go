package main

import (
	"github.com/hpcloud/log"
	"github.com/hpcloud/logyard-apps/applog_endpoint"
	"github.com/hpcloud/logyard-apps/applog_endpoint/config"
	"github.com/hpcloud/logyard-apps/applog_endpoint/drain"
)

func main() {
	config.LoadConfig()

	drain.RemoveOrphanedDrains()
	applog_endpoint.RouterMain()

	err := applog_endpoint.Serve()
	log.Fatal(err)
}
