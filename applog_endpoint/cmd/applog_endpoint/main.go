package main

import (
	"github.com/ActiveState/log"
	"github.com/ActiveState/logyard-apps/applog_endpoint"
	"github.com/ActiveState/logyard-apps/applog_endpoint/config"
	"github.com/ActiveState/logyard-apps/applog_endpoint/drain"
)

func main() {
	config.LoadConfig()

	drain.RemoveOrphanedDrains()
	applog_endpoint.RouterMain()

	err := applog_endpoint.Serve()
	log.Fatal(err)
}
