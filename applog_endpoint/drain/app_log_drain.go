package drain

import (
	"fmt"
	"logyard"
	"logyard/drain"
	"logyard/util/lineserver"
	"time"

	"github.com/hpcloud/log"
	"github.com/hpcloud/logyard-apps/applog_endpoint/config"
	"github.com/hpcloud/stackato-go/server"
)

const DRAIN_PREFIX = "tmp.applog_endpoint"

// UDP can lead to data loss, but that's OK for log streaming.
const LINESERVER_PROTO = "udp"

type AppLogDrain struct {
	appGUID   string
	srv       *lineserver.LineServer
	port      int
	lifetime  time.Duration // Keep the drain alive for this long
	drainName string
}

func NewAppLogDrain(appGUID string) (*AppLogDrain, error) {
	d := new(AppLogDrain)

	srv, err := lineserver.NewLineServer(
		LINESERVER_PROTO,
		fmt.Sprintf("%v:0", server.LocalIPMust()))
	if err != nil {
		return nil, err
	}

	addr, err := srv.GetUDPAddr()
	if err != nil {
		return nil, err
	}

	d.appGUID = appGUID
	d.srv = srv
	d.port = addr.Port
	d.lifetime, err = time.ParseDuration(config.GetConfig().DrainLifetime)
	if err != nil {
		return nil, fmt.Errorf(
			"Invalid duration value (%v) for drain_lifetime", err)
	}
	// TODO: name should have an uniq id, to allow multiple taile
	// sessions for same app.
	d.drainName = fmt.Sprintf("%s.%s", DRAIN_PREFIX, d.appGUID)

	return d, nil
}

func (d *AppLogDrain) Start() (chan string, error) {
	go d.srv.Start()

	err := d.addDrain()
	if err != nil {
		return nil, err
	}

	go func() {
		select {
		case <-time.After(d.lifetime):
			d.Stop(fmt.Errorf("Timed out %v", d.lifetime))
		case <-d.srv.Dying():
		}
	}()

	return d.srv.Ch, nil
}

func (d *AppLogDrain) Stop(reason error) {
	log.Infof("Stopping drain %s for reason: %v", d.Id(), reason)
	if err := d.removeDrain(); err != nil {
		log.Errorf("Failed to remove drain %v: %v", d.Id(), err)
	}
	d.srv.Kill(reason)
}

func (d *AppLogDrain) Wait() error {
	return d.srv.Wait()
}

func (d *AppLogDrain) Id() string {
	return fmt.Sprintf("%s<port:%d>", d.drainName, d.port)
}

// addDrain adds a logyard drain for the apptail.{appGUID} stream
// pointing to ourself (port)
func (d *AppLogDrain) addDrain() error {
	uri := fmt.Sprintf("%s://%v:%v",
		LINESERVER_PROTO,
		server.LocalIPMust(),
		d.port)
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
	log.Infof("Removed drain %v", d.Id())
	return err
}
