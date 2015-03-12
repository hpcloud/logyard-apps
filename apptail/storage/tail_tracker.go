package storage

import (
	"fmt"
	"github.com/ActiveState/log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type Tracker interface {
	Update(instKey string, childKey string, childVal int64)
	LoadTailers()
	Remove(key string)
	RegisterInstance(instKey string)
	InitializeChildNode(instKey string, childkey string, offSet int64)
	Submit()
	StartSubmissionTimer(retentionPeriod time.Duration)
	IsInstanceRegistered(instKey string) bool
	IsChildNodeInitialized(instKey string, childkey string) bool
	GetFileCachedOffset(instkey string, fname string) int64
}

type BoxedInt64 struct{ V int64 }

type TailNode map[string]*BoxedInt64

type Tailer struct {
	Instances map[string]TailNode
}

type tracker struct {
	storage       Storage
	Cached        *Tailer // do not expose this, it should ONLY be updated via Tracker methods
	mux           *sync.Mutex
	timerStopChan chan struct{} // used to send quit signal to timer
}

var (
	MinIOTicker = 5 * time.Second
)

func NewTracker(s Storage) Tracker {
	return &tracker{
		storage: s,
		mux:     &sync.Mutex{},
		Cached: &Tailer{
			Instances: make(map[string]TailNode),
		},
		timerStopChan: make(chan struct{}),
	}
}

func (t *tracker) StartSubmissionTimer(retentionPeriod time.Duration) {
	if retentionPeriod.Seconds() <= MinIOTicker.Seconds() {
		seconds := retentionPeriod / (1000 * time.Millisecond)
		log.Warnf("IMPORTANT: Setting retention period to %ds will increase your IO Rate", seconds)

	}
	ticker := time.NewTicker(retentionPeriod)
	go func() {
		for {
			select {
			case <-ticker.C:
				t.Submit()
			case <-t.timerStopChan:
				ticker.Stop()
				return
			}

		}
	}()
}

func (t *tracker) RegisterInstance(instKey string) {
	t.mux.Lock()
	if _, instance_exist := t.Cached.Instances[instKey]; !instance_exist {
		t.Cached.Instances[instKey] = TailNode{}
		t.formatMap("Current Status")
	}
	t.mux.Unlock()
}

// this is mainly used for testing since we are not exposing Cached via interface
func (t *tracker) IsInstanceRegistered(instKey string) bool {
	var exist bool
	t.mux.Lock()
	if _, instance_exist := t.Cached.Instances[instKey]; instance_exist {
		exist = instance_exist
	}
	t.mux.Unlock()
	return exist
}

func (t *tracker) IsChildNodeInitialized(instKey string, childkey string) bool {
	var exist bool
	t.mux.Lock()
	if tailNode, instance_exist := t.Cached.Instances[instKey]; instance_exist {
		if _, childNode_exist := tailNode[childkey]; childNode_exist {
			exist = childNode_exist
		}
	}
	t.mux.Unlock()
	return exist
}

func (t *tracker) InitializeChildNode(instKey string, childkey string, offSet int64) {
	t.mux.Lock()
	if tailNode, instance_exist := t.Cached.Instances[instKey]; instance_exist {
		if _, childNode_exist := tailNode[childkey]; !childNode_exist {
			tailNode[childkey] = &BoxedInt64{V: offSet}
			t.Cached.Instances[instKey] = tailNode
			t.formatMap("Current Status")
		}
	}
	t.mux.Unlock()
	runtime.Gosched()
}

func (t *tracker) GetFileCachedOffset(instkey string, fname string) int64 {
	var offset int64
	t.mux.Lock()
	if tailNode, instance_exist := t.Cached.Instances[instkey]; instance_exist {
		offset = atomic.LoadInt64(&tailNode[fname].V)
	}
	t.mux.Unlock()
	runtime.Gosched()
	return offset
}

func (t *tracker) Update(instKey string, childKey string, childVal int64) {
	if tailNode, instance_exist := t.Cached.Instances[instKey]; instance_exist {
		if _, childNode_exist := tailNode[childKey]; childNode_exist {
			atomic.StoreInt64(&tailNode[childKey].V, childVal)
		}
	}
}

func (t *tracker) Remove(key string) {
	log.Info("Removing the following key %s from cached instances", key)
	t.mux.Lock()
	delete(t.Cached.Instances, key)
	t.mux.Unlock()
	t.Submit()
}

func (t *tracker) LoadTailers() {
	t.mux.Lock()
	t.storage.Load(&t.Cached)
	t.formatMap("Loaded")
	t.mux.Unlock()
}

func (t *tracker) Submit() {
	t.mux.Lock()
	t.formatMap("Storing")
	t.storage.Write(t.Cached)
	t.mux.Unlock()
}

func (t *tracker) formatMap(ops string) {

	for k, v := range t.Cached.Instances {
		message := fmt.Sprintf("[%s] ContainerId: %s", ops, k)
		for fname, buffer := range v {

			log.Infof(message+" File: %s --> TailOffset: %d", fname, buffer)

		}

	}
}
