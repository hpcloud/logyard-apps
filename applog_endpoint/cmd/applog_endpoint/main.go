package main

import (
	"github.com/ActiveState/log"
	"github.com/ActiveState/logyard-apps/applog_endpoint"
	"github.com/ActiveState/logyard-apps/applog_endpoint/config"
)

func main() {
	// TODO: remove orphaned drains
	config.LoadConfig()

	applog_endpoint.RouterMain()

	err := applog_endpoint.Serve()
	log.Fatal(err)
}
