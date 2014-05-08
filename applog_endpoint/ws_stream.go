package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	"github.com/ActiveState/log"
)

// WebSocketStream wraps a websocket connection to provide error handling
type WebSocketStream struct {
	*websocket.Conn
}

type wsStreamData struct {
	Err   string `json:"error"`
	Value string `json:"value"`
}

// Send sends the value back to the client
func (s *WebSocketStream) Send(value string) error {
	return s.send(&wsStreamData{"", value})
}

// Fatalf sends the error back to the client, and closes the connection
func (s *WebSocketStream) Fatalf(format string, v ...interface{}) {
	err := s.send(&wsStreamData{fmt.Sprintf(format, v...), ""})
	if err != nil {
		log.Warnf("Error sending error back to websocket client: %v", err)
	}
	s.Close()
}

func (s *WebSocketStream) Fatal(v ...interface{}) {
	s.Fatal(fmt.Sprint(v...))
}

func (s *WebSocketStream) send(data *wsStreamData) error {
	jdata, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = s.Write(jdata)
	return err
}
