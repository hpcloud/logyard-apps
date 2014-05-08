package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/ActiveState/log"
	"io"
	"net/http"
)

func echoHandler(ws *websocket.Conn) {
	io.Copy(ws, ws)
}

func readArguments(ws *websocket.Conn) (token, appGUID string, err error) {
	q := ws.Config().Location.Query()
	appGUID = q.Get("appid")
	token = ws.Config().Header.Get("Authorization")
	if token == "" {
		token = q.Get("token")
	}
	if token == "" {
		err = fmt.Errorf("empty token")
	} else if appGUID == "" {
		err = fmt.Errorf("missing appGUID")
	}
	return
}

func tailHandler(ws *websocket.Conn) {
	log.Infof("tailHandler start %+v", ws)
	stream := &WebSocketStream{ws}
	token, appGUID, err := readArguments(ws)
	if err != nil {
		stream.Fatal(err)
		return
	}

	// First authorize with the CC by fetching something
	_, err = recentLogs(token, appGUID, 1)
	if err != nil {
		stream.Fatal(err)
		return
	}

	drain, err := NewAppLogDrain(appGUID)
	if err != nil {
		stream.Fatal(err)
		return
	}
	ch, err := drain.Start()
	if err != nil {
		stream.Fatal(err)
	}
	for line := range ch {
		stream.Send(line)
	}
	log.Infof("tailHandler done")
}

func serve() error {
	addr := fmt.Sprintf(":%d", PORT)
	http.Handle("/echo", websocket.Handler(echoHandler))
	http.Handle("/tail", websocket.Handler(tailHandler))
	return http.ListenAndServe(addr, nil)
}
