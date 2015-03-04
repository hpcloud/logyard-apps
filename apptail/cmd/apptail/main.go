package main

import (
	"github.com/ActiveState/log"
	"github.com/ActiveState/logyard-apps/apptail"
	"github.com/ActiveState/logyard-apps/apptail/docker"
	apptail_event "github.com/ActiveState/logyard-apps/apptail/event"
	"github.com/ActiveState/logyard-apps/common"
	"github.com/ActiveState/stackato-go/server"
	"github.com/alecthomas/gozmq"
	"github.com/apcera/nats"
	uuid "github.com/nu7hatch/gouuid"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
	"time"
)

func main() {
	go common.RegisterTailCleanup()
	major, minor, patch := gozmq.Version()
	log.Infof("Starting apptail (zeromq %d.%d.%d)", major, minor, patch)

	apptail.LoadConfig()
	log.Infof("Config: %+v\n", apptail.GetConfig())

	uid := getUID()

	natsclient := server.NewNatsClient(3)

	var state = make(map[string]string)
	var mutex = sync.Mutex{}

	var key_exist bool

	mutex.Lock()
	if _, key_exist = state[uid]; !key_exist {
		state[uid] = time.Now().Local().Format(time.RFC3339)
	}
	log.Infof("registered instances with nats: %s", state)
	mutex.Unlock()
	runtime.Gosched()

	if !key_exist {

		natsclient.Subscribe("logyard."+uid+".newinstance", func(instance *apptail.Instance) {
			instance.Tail()
		})

		natsclient.Publish("logyard."+uid+".start", []byte("{}"))

		log.Infof("Waiting for app instances ...")

		go docker.DockerListener.Listen()

		server.MarkRunning("apptail")

		apptail_event.MonitorCloudEvents()
	} else {
		log.Infof("already subscribed to nats :", uid)

	}
}

func track(uid string, state map[string]string, mux sync.Mutex, natsclient *nats.EncodedConn) (wait <-chan struct{}) {
	instance_name := "logyard." + uid + ".newinstance"

	ch := make(chan struct{})
	mux.Lock()
	_, Key_exist := state[uid]
	mux.Unlock()
	runtime.Gosched()
	if !Key_exist {

		natsclient.Subscribe(instance_name, func(instance *apptail.Instance) {
			instance.Tail()
		})

		natsclient.Publish("logyard."+uid+".start", []byte("{}"))

	} else {
		close(ch)

	}

	return ch

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
