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
	LoadTailers()
	Remove(key string)
	Status()
	RegisterInstance(instKey string)
	InitializeChildNode(instKey string, childkey string, offSet int64)
	Submit()
}

type TailNode map[string]int64

type Tailer struct {
	IsLive    bool
	Instances map[string]TailNode // maybe a pointer
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

func(t *tracker) Status(){
	t.mux.Lock()
	log.Info("---DEBUG---", t.Cached.Instances)

	t.mux.Unlock()
	runtime.Gosched()
}

func (t *tracker) RegisterInstance(instKey string){
	t.mux.Lock()
	if _, instance_exist := t.Cached.Instances[instKey]; !instance_exist{
		t.Cached.IsLive = true
		t.Cached.Instances[instKey] = TailNode{}
		log.Info("[RegisterInstance]Current status : ", t.Cached.Instances)
	}
	t.mux.Unlock()
	runtime.Gosched()
	
}

func(t *tracker) InitializeChildNode(instKey string, childkey string, offSet int64){
	t.mux.Lock()
	if _, instance_exist := t.Cached.Instances[instKey]; instance_exist{
		tailNode := t.Cached.Instances[instKey]
		if _, childNode_exist := tailNode[childkey]; !childNode_exist{
			tailNode[childkey] = offSet
			t.Cached.Instances[instKey] = tailNode
			log.Info("[InitializeChildNode]Current status : ", t.Cached.Instances)
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
			//log.Info("[Update]Current status : ", t.Cached.Instances)
		}
		
	}
	t.mux.Unlock()
	runtime.Gosched()
}

func (t *tracker) LoadTailers(){
	t.mux.Lock()
	t.storage.Load(&t.Cached)
	log.Info("[LoadTailers]Loaded the following tailers from previous session:", t.Cached.Instances)
	t.mux.Unlock()
	runtime.Gosched()
}

func(t *tracker) Submit(){
	t.mux.Lock()
	log.Info("Storing the offset in the following instances:", t.Cached.Instances)
	t.storage.Write(t.Cached)
	log.Info("[Submit]Current status : ", t.Cached.Instances)
	t.mux.Unlock()
	runtime.Gosched()
}

func(t *tracker) Remove(key string) {
	log.Info("Removing the following key:", key)

	t.mux.Lock()
	delete(t.Cached.Instances, key)
	log.Info("[Remove]Current status : ", t.Cached.Instances)
	t.mux.Unlock()
	runtime.Gosched()

}
