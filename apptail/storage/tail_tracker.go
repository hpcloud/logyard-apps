package storage

import (
	"fmt"
	"github.com/ActiveState/log"
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
	Commit() error
	StartSubmissionTimer(persistInterval time.Duration)
	IsInstanceRegistered(instKey string) bool
	IsChildNodeInitialized(instKey string, childkey string) bool
	GetFileCachedOffset(instkey string, fname string) int64
	getBuffer() ([]byte, error)
	CleanUp(validIds map[string]bool)
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
	writeMux      *sync.Mutex
	timerStopChan chan struct{} // used to send quit signal to time,
	debug         bool
}

func NewTracker(s Storage, debug bool) Tracker {

	return &tracker{
		storage:  s,
		mux:      &sync.Mutex{},
		writeMux: &sync.Mutex{},
		Cached: &Tailer{
			Instances: make(map[string]TailNode),
		},
		timerStopChan: make(chan struct{}),
		debug:         debug,
	}
}

func (t *tracker) StartSubmissionTimer(persistInterval time.Duration) {
	ticker := time.NewTicker(persistInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				err := t.Commit()
				if err != nil {
					log.Fatal(err)

				}
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

		t.dumpState("Current Status")

	}
	t.mux.Unlock()
}

// this is mainly used for testing since we are not exposing Cached via interface
func (t *tracker) IsInstanceRegistered(instKey string) bool {
	t.mux.Lock()
	defer t.mux.Unlock()
	_, exist := t.Cached.Instances[instKey]

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
			t.dumpState("Current Status")

		}
	}
	t.mux.Unlock()
}

func (t *tracker) GetFileCachedOffset(instkey string, fname string) int64 {
	var offset int64
	t.mux.Lock()
	if tailNode, instance_exist := t.Cached.Instances[instkey]; instance_exist {
		offset = atomic.LoadInt64(&tailNode[fname].V)
	}
	t.mux.Unlock()
	return offset
}

func (t *tracker) Update(instKey string, childKey string, childVal int64) {
	if tailNode, instance_exist := t.Cached.Instances[instKey]; instance_exist {
		if _, childNode_exist := tailNode[childKey]; childNode_exist {
			atomic.StoreInt64(&tailNode[childKey].V, childVal)
		}
	}
}

func (t *tracker) CleanUp(validIds map[string]bool) {
	cachedIds := make(map[string]bool)
	t.mux.Lock()
	for key := range t.Cached.Instances {
		cachedIds[key] = true

	}
	instancesToRemove := getInvalidInstances(cachedIds, validIds)
	for key := range instancesToRemove {
		delete(t.Cached.Instances, key)

	}

	t.dumpState("Cleaned up")

	t.mux.Unlock()

}

func (t *tracker) Remove(key string) {
	if t.debug {
		log.Infof("Removing the following key %s from cached instances", key)
	}
	t.mux.Lock()
	delete(t.Cached.Instances, key)
	t.mux.Unlock()
	if err := t.Commit(); err != nil {
		log.Fatal(err)

	}
}

func (t *tracker) LoadTailers() {
	t.mux.Lock()
	t.storage.Load(&t.Cached)
	t.dumpState("Loaded")
	t.mux.Unlock()
}

func (t *tracker) getBuffer() ([]byte, error) {
	t.mux.Lock()
	defer t.mux.Unlock()
	return t.storage.Encode(t.Cached)

}

func (t *tracker) Commit() error {

	bytes, err := t.getBuffer()

	if err != nil {
		return err

	}

	//t.dumpState("Storing")

	t.writeMux.Lock()
	err = t.storage.Write(bytes)
	if err != nil {
		return err

	}
	t.writeMux.Unlock()
	return nil
}

func (t *tracker) dumpState(ops string) {
	// Important to wrap this in a mutex since it's accessing shared resource
	if t.debug {
		for k, v := range t.Cached.Instances {
			message := fmt.Sprintf("[%s] ContainerId: %s", ops, k)
			for fname, buffer := range v {

				log.Infof(message+" File: %s --> TailOffset: %d", fname, buffer)

			}

		}
	}
}

func getInvalidInstances(dockerIds map[string]bool, allInstances map[string]bool) map[string]bool {
	invalidInstances := make(map[string]bool)
	for id, _ := range dockerIds {
		if _, exist := allInstances[id]; !exist {
			invalidInstances[id] = true
		}
	}
	return invalidInstances

}
