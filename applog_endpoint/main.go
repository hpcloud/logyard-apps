package main

import (
	"github.com/ActiveState/log"
)

const COMPONENT = "websocket_endpoint"
const PORT = 5722

func main() {
	if err := advertiseToStackatoRouter(); err != nil {
		log.Fatal(err)
	}

	log.Fatal(serve())
}
