package main

import (
	"flag"
	"fmt"
	"github.com/ActiveState/log"
	"github.com/ActiveState/logyard-apps/apptail"
	"github.com/ActiveState/logyard-apps/apptail/docker"
	apptail_event "github.com/ActiveState/logyard-apps/apptail/event"
	"github.com/ActiveState/logyard-apps/apptail/storage"
	"github.com/ActiveState/logyard-apps/common"
	"github.com/ActiveState/stackato-go/server"
	"github.com/alecthomas/gozmq"
	uuid "github.com/nu7hatch/gouuid"
	"io/ioutil"
	"os"
	"time"
)

var (
	stateFile_path = flag.String("state_file_path", default_path, "Path to the state file")

	default_path = fmt.Sprintf("%s/.apptail.gob", os.Getenv("HOME"))

	debug = flag.Bool("debug", false, "debugger")
)

func main() {
	flag.Parse()
	go common.RegisterTailCleanup()

	apptail.LoadConfig()

	fstorage := storage.NewFileStorage(*stateFile_path)
	tracker := storage.NewTracker(fstorage, *debug)
	tracker.LoadTailers()

	interval := time.Duration(int64(apptail.GetConfig().PersistPositionIntervalSeconds))
	go tracker.StartSubmissionTimer(interval * time.Second)

	major, minor, patch := gozmq.Version()
	log.Infof("Starting apptail (zeromq %d.%d.%d)", major, minor, patch)

	log.Infof("Config: %+v\n", apptail.GetConfig())

	uid := getUID()

	natsclient := server.NewNatsClient(3)

	natsclient.Subscribe("logyard."+uid+".newinstance", func(instance *apptail.Instance) {
		instance.Tail(tracker)
	})

	natsclient.Publish("logyard."+uid+".start", []byte("{}"))
	log.Infof("Waiting for app instances ...")

	go docker.DockerListener.Listen()

	// clean up the cache after restart
	docker.DockerListener.TrackerCleanUp(tracker)

	server.MarkRunning("apptail")

	apptail_event.MonitorCloudEvents()

}

// getUID returns the UID of the aggregator running on this node. the UID is
// also shared between the local dea/stager, so that we send/receive messages
// only from the local dea/stagers.
func getUID() string {
	var UID string
	uidFile := "/tmp/logyard.uid"
	if _, err := os.Stat(uidFile); os.IsNotExist(err) {
		uid, err := uuid.NewV4()
		if err != nil {
			common.Fatal("%v", err)
		}
		UID = uid.String()
		if err = ioutil.WriteFile(uidFile, []byte(UID), 0644); err != nil {
			common.Fatal("%v", err)
		}
	} else {
		data, err := ioutil.ReadFile(uidFile)
		if err != nil {
			common.Fatal("%v", err)
		}
		UID = string(data)
	}
	log.Infof("detected logyard UID: %s\n", UID)
	return UID
}
