package wsutil

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

type Handler interface {
	ServeWS(http.ResponseWriter, *http.Request, *WebSocketStream)
}

type HandlerFunc func(http.ResponseWriter, *http.Request, *WebSocketStream)

func (f HandlerFunc) ServeWS(w http.ResponseWriter, r *http.Request, ws *WebSocketStream) {
	f(w, r, ws)
}

type webSocketHandler struct {
	handler Handler
}

func WebSocketHandler(h Handler) http.Handler {
	return &webSocketHandler{h}
}

func (h *webSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		var errString string
		if _, ok := err.(websocket.HandshakeError); !ok {
			errString = fmt.Sprintf("Handshake error: %v", err)
		} else {
			errString = fmt.Sprintf("Unknown websocket error: %v", err)
		}
		log.Info(errString)
		http.Error(w, errString, 500)
		return
	}

	log.Infof("wsutil.ServeWS start - %v", getWsConnId(r, ws))
	defer log.Infof("wsutil.ServeWS finish - %v", getWsConnId(r, ws))

	h.handler.ServeWS(w, r, &WebSocketStream{ws})

	ws.Close()
}

// XXX: pass this as a log context (gorilla) object
func getWsConnId(r *http.Request, ws *websocket.Conn) string {
	return fmt.Sprintf("ws:/%v %v (subprotocol %+v)",
		r.URL.Path, ws.RemoteAddr(), ws.Subprotocol())
}
