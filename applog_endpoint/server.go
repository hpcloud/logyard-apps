package main

import (
	"fmt"
	"github.com/ActiveState/log"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Ignore origin checks (won't work with wscat)
		return true
	},
}

func readArguments(r *http.Request) (token, appGUID string, err error) {
	appGUID = r.FormValue("appid")
	token = r.Header.Get("Authorization")
	if token == "" {
		token = r.FormValue("token")
	}
	if token == "" {
		err = fmt.Errorf("empty token")
	} else if appGUID == "" {
		err = fmt.Errorf("missing appGUID")
	}
	return
}

func getWsConnId(r *http.Request, ws *websocket.Conn) string {
	return fmt.Sprintf("ws:/%v %v (subprotocol %+v)",
		r.URL.Path, ws.RemoteAddr(), ws.Subprotocol())
}

func tailHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Infof("Handshake error: %v", err)
		} else {
			log.Infof("Unknown websocket error: %v", err)
		}
		return
	}

	token, appGUID, err := readArguments(r)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Invalid arguments; %v", err), 400)
		return
	}

	tailHandlerWs(r, ws, token, appGUID)
}

func tailHandlerWs(
	r *http.Request, ws *websocket.Conn, token, appGUID string) {
	log.Infof("WS init - %v", getWsConnId(r, ws))
	defer log.Infof("WS done - %v", getWsConnId(r, ws))

	stream := &WebSocketStream{ws}

	// First authorize with the CC by fetching something
	_, err := recentLogs(token, appGUID, 1)
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

	// TODO: don't block here forever; handle client disconnections.
	// else, we keep the drain open for 20m.
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
	http.HandleFunc("/tail", tailHandler)
	return http.ListenAndServe(addr, nil)
}
