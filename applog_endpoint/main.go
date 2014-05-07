package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/ActiveState/log"
	"io"
	"net/http"
)

func echoHandler(ws *websocket.Conn) {
	io.Copy(ws, ws)
}

func main() {
	if err := advertiseToStackatoRouter(); err != nil {
		log.Fatal(err)
	}

	http.Handle("/echo", websocket.Handler(echoHandler))
	err := http.ListenAndServe(":5722", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
