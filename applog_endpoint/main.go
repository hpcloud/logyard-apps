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

func logsHandler(ws *websocket.Conn) {
	q := ws.Config().Location.Query()
	appGUID := q.Get("appGUID")
	token := ws.Config().Header.Get("Authorization")
	if token == "" {
		token = q.Get("token")
	}
	if token == "" {
		io.WriteString(ws, "ERROR: empty token")
	} else if appGUID == "" {
		io.WriteString(ws, "ERROR: missing appGUID")
	} else {
		io.WriteString(ws, recentLogs(token, appGUID))
	}
}

func main() {
	if err := advertiseToStackatoRouter(); err != nil {
		log.Fatal(err)
	}

	http.Handle("/echo", websocket.Handler(echoHandler))
	http.Handle("/tail", websocket.Handler(logsHandler))
	err := http.ListenAndServe(":5722", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
