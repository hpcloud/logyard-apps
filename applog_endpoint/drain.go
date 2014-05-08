package main

import (
	"bufio"
	"fmt"
	"github.com/ActiveState/log"
	"logyard"
	"logyard/drain"
	"net"
	"stackato/server"
)

// addDrain adds a logyard drain for the apptail.{appGUID} stream
// pointing to ourself (port)
func addDrain(appGUID string, port int) error {
	name := fmt.Sprintf("tmp.websocket_endpoint.%s", appGUID)
	uri := fmt.Sprintf("udp://%v:%v", server.LocalIPMust(), port)
	filter := fmt.Sprintf("apptail.%s", appGUID)
	drainURI, err := drain.ConstructDrainURI(
		name, uri, []string{filter}, nil)
	if err != nil {
		return err
	}
	if err = logyard.AddDrain(name, drainURI); err != nil {
		return err
	}
	log.Infof("Added drain %v => %v", name, drainURI)
	return nil
}

func listenOnAppLogStream(appGUID string) (chan []byte, error) {
	// TODO: cleanup upon request and/or timeout
	globAddr, err := net.ResolveUDPAddr(
		"udp", fmt.Sprintf("%v:0", server.LocalIPMust()))
	if err != nil {
		return nil, err
	}
	sock, err := net.ListenUDP("udp", globAddr)
	if err != nil {
		return nil, fmt.Errorf("can't listen: %v", err)
	}

	addr, err := net.ResolveUDPAddr("udp", sock.LocalAddr().String())
	if err != nil {
		return nil, fmt.Errorf("failed to resolve udp addr: %v", err)
	}
	err = addDrain(appGUID, addr.Port)
	if err != nil {
		return nil, err
	}

	ch := make(chan []byte)
	reader := bufio.NewReader(sock)
	go func() {
		for {
			line, _, err := reader.ReadLine()
			if err != nil {
				ch <- []byte(fmt.Sprintf("INTERNAL ERROR %v", err))
				return
			}
			ch <- line
		}
	}()
	return ch, nil
}
