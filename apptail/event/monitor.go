package event

import (
	"encoding/json"
	"fmt"
	"github.com/ActiveState/log"
	"github.com/ActiveState/logyard-apps/apptail/message"
	"github.com/ActiveState/logyard-apps/apptail/util"
	"github.com/ActiveState/logyard-apps/common"
	"github.com/ActiveState/logyard-apps/sieve"
	"github.com/ActiveState/zmqpubsub"
	"logyard"
	"time"
)

// Make relevant cloud events available in application logs. Heroku style.
func MonitorCloudEvents() {
	sub := logyard.Broker.Subscribe("event.timeline")
	defer sub.Stop()

	pub := logyard.Broker.NewPublisherMust()
	defer pub.Stop()

	log.Info("Listening for app relevant cloud events...")
	for msg := range sub.Ch {
		var event sieve.Event

		err := json.Unmarshal([]byte(msg.Value), &event)
		if err != nil {
			common.Fatal("%v", err) // not expected at all
		}

		// Re-parse the event json record into a TimelineEvent structure.
		var t TimelineEvent
		if data, err := json.Marshal(event.Info); err != nil {
			common.Fatal("%v", err)
		} else {
			err = json.Unmarshal(data, &t)
			if err != nil {
				common.Fatal("Invalid timeline event: %v", err)
			}
		}

		var source string
		brandname := util.GetBrandName()
		if t.InstanceIndex > -1 {
			source = fmt.Sprintf("%v[%v.%v]", brandname, event.Process, t.InstanceIndex)
		} else {
			source = fmt.Sprintf("%v[%v]", brandname, event.Process)
		}

		PublishAppLog(pub, t, source, &event)
	}
	log.Warn("Finished listening for app relevant cloud events.")

	err := sub.Wait()
	if err != nil {
		common.Fatal("%v", err)
	}
}

func PublishAppLog(
	pub *zmqpubsub.Publisher,
	t TimelineEvent,
	source string, event *sieve.Event) {

	err := (&message.Message{
		LogFilename:   "",
		Source:        source,
		InstanceIndex: t.InstanceIndex,
		AppGUID:       t.App.GUID,
		AppName:       t.App.Name,
		AppSpace:      t.App.Space,
		MessageCommon: common.NewMessageCommon(event.Desc, time.Unix(event.UnixTime, 0), util.LocalNodeId()),
	}).Publish(pub, true)

	if err != nil {
		common.Fatal("%v", err)
	}
}
