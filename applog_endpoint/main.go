package main

import (
	"github.com/ActiveState/log"
)

const COMPONENT = "websocket_endpoint"
const PORT = 5722

func main() {
	// TODO: remove orphaned drains
	LoadConfig()

	routerMain()

	err := serve()

	log.Fatal(err)
}
