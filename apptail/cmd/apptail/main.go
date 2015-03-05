package main

import (
	"github.com/ActiveState/log"
	"github.com/ActiveState/logyard-apps/apptail"
	"github.com/ActiveState/logyard-apps/apptail/docker"
	apptail_event "github.com/ActiveState/logyard-apps/apptail/event"
	"github.com/ActiveState/logyard-apps/common"
	"github.com/ActiveState/stackato-go/server"
	"github.com/alecthomas/gozmq"
	uuid "github.com/nu7hatch/gouuid"
	"io/ioutil"
	"os"
	"sync"
)

func main() {
	go common.RegisterTailCleanup()
	major, minor, patch := gozmq.Version()
	log.Infof("Starting apptail (zeromq %d.%d.%d)", major, minor, patch)

	apptail.LoadConfig()
	log.Infof("Config: %+v\n", apptail.GetConfig())

	uid := getUID()

	natsclient := server.NewNatsClient(3)

	mux := &sync.Mutex{}

	n := 0
	started_instances := make(map[string]int)

	natsclient.Subscribe("logyard."+uid+".newinstance", func(instance *apptail.Instance) {
		n++
		mux.Lock()
		_, key_exist := started_instances[instance.DockerId]
		mux.Unlock()

		if !key_exist {
			mux.Lock()
			started_instances[instance.DockerId] = n
			mux.Unlock()
			go func() {
				instance.Tail()
				mux.Lock()
				delete(started_instances, instance.DockerId)
				log.Info("available instances: ", started_instances)
				mux.Unlock()
			}()
		}
	})

	natsclient.Publish("logyard."+uid+".start", []byte("{}"))
	log.Infof("Waiting for app instances ...")

	go docker.DockerListener.Listen()

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
