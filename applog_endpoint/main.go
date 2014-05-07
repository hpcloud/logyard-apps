package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	"github.com/ActiveState/log"
	"io"
	"net/http"
)

func echoHandler(ws *websocket.Conn) {
	io.Copy(ws, ws)
}

func logsHandler(ws *websocket.Conn) {
	q := ws.Config().Location.Query()
	appGUID := q.Get("appid")
	token := ws.Config().Header.Get("Authorization")
	if token == "" {
		token = q.Get("token")
	}
	// TODO: marshall errors in json
	if token == "" {
		io.WriteString(ws, "ERROR: empty token")
	} else if appGUID == "" {
		io.WriteString(ws, "ERROR: missing appGUID")
	} else {
		if logs, err := recentLogs(token, appGUID); err != nil {
			io.WriteString(ws, fmt.Sprintf("ERROR: %v", err))
		} else {
			for _, line := range logs {
				if output, err := json.Marshal(line); err != nil {
					io.WriteString(ws, fmt.Sprintf("ERROR: %v", err))
				} else {
					ws.Write(output)
				}
			}
		}
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
