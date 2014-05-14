package main

import (
	"fmt"
	"github.com/ActiveState/log"
	"github.com/gorilla/mux"
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

// XXX: pass this as a log context (gorilla) object
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

	args, err := ParseArguments(r)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Invalid arguments; %v", err), 400)
		return
	}

	tailHandlerWs(r, ws, args)
}

func tailHandlerWs(
	r *http.Request, ws *websocket.Conn, args *Arguments) {
	log.Infof("WS init - %v", getWsConnId(r, ws))
	defer log.Infof("WS done - %v", getWsConnId(r, ws))

	stream := &WebSocketStream{ws}

	// First authorize with the CC by fetching something
	_, err := recentLogs(args.Token, args.GUID, 1)
	if err != nil {
		stream.Fatalf("%v", err)
		return
	}

	// Recent history requested?
	if args.Num > 0 {
		recentLogs, err := recentLogs(args.Token, args.GUID, args.Num)
		if err != nil {
			stream.Fatalf("%v", err)
			return
		}
		for _, _ = range recentLogs {
			// TODO: make CC return raw log data from apptail as-is
			// without mangling.
			stream.Send("{}")
		}
	}

	drain, err := NewAppLogDrain(args.GUID)
	if err != nil {
		stream.Fatalf("Unable to create drain: %v", err)
		return
	}
	ch, err := drain.Start()
	if err != nil {
		stream.Fatalf("Unable to start drain: %v", err)
	}

	err = stream.Forward(ch)
	if err != nil {
		log.Infof("%v", err)
		drain.Stop(err)
	}

	// We expect drain.Wait to not block at this point.
	if err := drain.Wait(); err != nil {
		if _, ok := err.(WebSocketStreamError); !ok {
			log.Warnf("Error from app log drain server: %v", err)
		}
	}
}

func serve() error {
	addr := fmt.Sprintf(":%d", PORT)
	r := mux.NewRouter()
	r.HandleFunc("/v2/apps/{guid}/tail", tailHandler)
	// r.HandleFunc("/v2/apps/{guid}/recent", recentHandler)
	http.Handle("/", r)
	return http.ListenAndServe(addr, nil)
}
