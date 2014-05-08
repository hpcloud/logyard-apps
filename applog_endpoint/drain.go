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
	appGUID   string
	srv       *lineserver.LineServer
	port      int
	drainName string
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
	// TODO: name should have an uniq id, to allow multiple taile
	// sessions for same app.
	d.drainName = fmt.Sprintf("tmp.websocket_endpoint.%s", d.appGUID)

	return d, nil
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
	d.srv.Kill(d.removeDrain())
}

// addDrain adds a logyard drain for the apptail.{appGUID} stream
// pointing to ourself (port)
func (d *AppLogDrain) addDrain() error {
	uri := fmt.Sprintf("udp://%v:%v", server.LocalIPMust(), d.port)
	filter := fmt.Sprintf("apptail.%s", d.appGUID)
	drainURI, err := drain.ConstructDrainURI(
		d.drainName, uri, []string{filter}, nil)
	if err != nil {
		return err
	}
	if err = logyard.AddDrain(d.drainName, drainURI); err != nil {
		return err
	}
	log.Infof("Added drain %v => %v", d.drainName, drainURI)
	return nil
}

func (d *AppLogDrain) removeDrain() error {
	err := logyard.DeleteDrain(d.drainName)
	log.Infof("Removed drain %v", d.drainName)
	return err
}
