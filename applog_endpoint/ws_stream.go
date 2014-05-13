package main

import (
	"fmt"
	"github.com/ActiveState/log"
	"github.com/gorilla/websocket"
	"time"
)

// WebSocketStream wraps a websocket connection to provide error handling
type WebSocketStream struct {
	*websocket.Conn
}

type WebSocketStreamError struct {
	error
}

type wsStreamData struct {
	Err   string      `json:"error"`
	Value interface{} `json:"value"`
}

// Forward reads from channel and sends the values. Also pings the
// client periodically.
func (s *WebSocketStream) Forward(ch chan string) error {
	for {
		select {
		case line, ok := <-ch:
			if !ok {
				return nil // All done.
			}
			if err := s.Send(line); err != nil {
				return WebSocketStreamError{fmt.Errorf(
					"Closing websocket because of write error: %v", err)}
			}
		case <-time.After(time.Second):
			// Check if client is alive every second
			err := s.WriteControl(
				websocket.PingMessage, nil, time.Now().Add(time.Second))
			if err != nil {
				return WebSocketStreamError{fmt.Errorf(
					"Closing websocket because of ping error: %v", err)}
			}

		}
	}
}

// Send sends the value back to the client
func (s *WebSocketStream) Send(value string) error {
	return s.send(&wsStreamData{"", value})
}

// Fatalf sends the error back to the client, and closes the connection
func (s *WebSocketStream) Fatalf(format string, v ...interface{}) {
	data := &wsStreamData{fmt.Sprintf(format, v...), ""}
	err := s.send(data)
	if err != nil {
		log.Warnf("Error sending error back to websocket client: %v", err)
	}
	s.Close()
}

func (s *WebSocketStream) send(data *wsStreamData) error {
	return s.WriteJSON(data)
}
