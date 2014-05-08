package main

import (
	"fmt"
	"github.com/ActiveState/log"
	"logyard"
	"logyard/drain"
	"logyard/util/lineserver"
	"stackato/server"
)

type AppLogDrain struct {
	appGUID string
	srv     *lineserver.LineServer
	port    int
}

func NewAppLogDrain(appGUID string) (*AppLogDrain, error) {
	d := new(AppLogDrain)

	srv, err := lineserver.NewLineServer(
		fmt.Sprintf("%v:0", server.LocalIPMust()))
	if err != nil {
		return nil, err
	}

	addr, err := srv.GetAddr()
	if err != nil {
		return nil, err
	}

	d.appGUID = appGUID
	d.srv = srv
	d.port = addr.Port

	return d, nil
}

// addDrain adds a logyard drain for the apptail.{appGUID} stream
// pointing to ourself (port)
func (d *AppLogDrain) addDrain() error {
	// TODO: name should have an uniq id, to allow multiple taile
	// sessions for same app.
	name := fmt.Sprintf("tmp.websocket_endpoint.%s", d.appGUID)
	uri := fmt.Sprintf("udp://%v:%v", server.LocalIPMust(), d.port)
	filter := fmt.Sprintf("apptail.%s", d.appGUID)
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

func (d *AppLogDrain) Start() (chan string, error) {
	// TODO: cleanup upon request and/or timeout
	go d.srv.Start()

	err := d.addDrain()
	if err != nil {
		return nil, err
	}

	return d.srv.Ch, nil
}

func (d *AppLogDrain) Stop() {
	d.srv.Kill(nil)
}
