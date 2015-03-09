package storage

import (
	"runtime"
	"fmt"
	"github.com/ActiveState/log"
	"os"
	"sync"
)

type Tracker interface {
	Update(key string, childKey string, childVal int64)
	LoadTailers() *Tailer
	Remove(key string)
	Status() map[string]TailNode
	RegisterInstance(instKey string)
	InitializeChildNode(instKey string, childkey string)
}

type TailNode map[string]int64

type Tailer struct {
	IsLive    bool
	Instances map[string]TailNode
}

type tracker struct {
	storage Storage
	Cached  *Tailer
	mux     *sync.Mutex
}

var (
	Path = fmt.Sprintf("%s/.apptail.gob", os.Getenv("HOME"))
)

func NewTracker(s Storage) Tracker {
	return &tracker{
		storage: s,
		mux:     &sync.Mutex{},
		Cached: &Tailer{
			Instances: make(map[string]TailNode), // we only need to instanciate this once
		},
	}
}

func(t *tracker) Status() map[string]TailNode{
	return t.Cached.Instances
	
}

func (t *tracker) LoadTailers() *Tailer {
	t.mux.Lock()
	t.storage.Load(&t.Cached)
	log.Info("loaded the following tailers from previous session:", t.Cached.Instances)
	t.mux.Unlock()
	runtime.Gosched()
	return t.Cached
}


func (t *tracker) RegisterInstance(instKey string){
	t.mux.Lock()
	if _, instance_exist := t.Cached.Instances[instKey]; !instance_exist{
		t.Cached.IsLive = true
		t.Cached.Instances[instKey] = TailNode{}
		
	}
	t.mux.Unlock()
	runtime.Gosched()
	
}

func(t *tracker) InitializeChildNode(instKey string, childkey string){
	t.mux.Lock()
	if _, instance_exist := t.Cached.Instances[instKey]; instance_exist{
		tailNode := t.Cached.Instances[instKey]
		if _, childNode_exist := tailNode[childkey]; !childNode_exist{
			tailNode[childkey] = 0
			t.Cached.Instances[instKey] = tailNode
		}

	}
	t.mux.Unlock()
	runtime.Gosched()
}

func (t *tracker) Update(instKey string, childKey string, childVal int64){
	t.mux.Lock()
	if _, instance_exist := t.Cached.Instances[instKey]; instance_exist{
		tailNode := t.Cached.Instances[instKey]
		if _, childNode_exist := tailNode[childKey]; childNode_exist{
			tailNode[childKey] = childVal
			t.Cached.Instances[instKey] = tailNode
		}
		
	}
	t.mux.Unlock()
	runtime.Gosched()
	
}


func (t *tracker) Remove(key string) {
	log.Info("Removing the following key:", key)
	t.mux.Lock()
	delete(t.Cached.Instances, key)
	t.mux.Unlock()
	runtime.Gosched()

}
