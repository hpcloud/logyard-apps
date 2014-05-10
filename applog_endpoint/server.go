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

func getWsConnId(ws *websocket.Conn) string {
	config := ws.Config()
	req := ws.Request()
	return fmt.Sprintf("ws:/%v %v (proto %+v; version %v)",
		req.URL.Path, req.RemoteAddr,
		config.Protocol, config.Version)
}

func tailHandler(ws *websocket.Conn) {
	log.Infof("WS init - %v", getWsConnId(ws))
	defer log.Infof("WS done - %v", getWsConnId(ws))

	stream := &WebSocketStream{ws}
	token, appGUID, err := readArguments(ws)
	if err != nil {
		stream.Fatalf("Invalid arguments: %v", err)
		return
	}

	// First authorize with the CC by fetching something
	_, err = recentLogs(token, appGUID, 1)
	if err != nil {
		stream.Fatalf("%v", err)
		return
	}

	drain, err := NewAppLogDrain(appGUID)
	if err != nil {
		stream.Fatalf("Unable to create drain: %v", err)
		return
	}
	ch, err := drain.Start()
	if err != nil {
		stream.Fatalf("Unable to start drain: %v", err)
	}

	for line := range ch {
		if err := stream.Send(line); err != nil {
			log.Infof("Closing websocket because of write error: %v", err)
			drain.Stop(err)
			return
		}
	}

	if err := drain.Wait(); err != nil {
		log.Warnf("Error from app log drain server: %v", err)
	}

}

func serve() error {
	addr := fmt.Sprintf(":%d", PORT)
	http.Handle("/echo", websocket.Handler(echoHandler))
	http.Handle("/tail", websocket.Handler(tailHandler))
	return http.ListenAndServe(addr, nil)
}
