package storage

import (
	"runtime"
	"fmt"
	"github.com/ActiveState/log"
	"os"
	"sync"
)

type Tracker interface {
	Update(instKey string, childKey string, childVal int64)
	LoadTailers()
	Remove(key string)
	RegisterInstance(instKey string)
	InitializeChildNode(instKey string, childkey string, offSet int64)
	Submit()
	GetRemoveChan()chan string
	GetCommitChan()chan bool
	GetUpdateChan()chan bool
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
	removeChan chan string
	updateChan chan bool
	writeChan chan bool
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

func(t *tracker) GetRemoveChan()chan string{
	t.removeChan = make(chan string)
	return t.removeChan

}

func(t *tracker) GetCommitChan()chan bool{
	t.writeChan = make(chan bool)
	return t.writeChan
	
}

func(t *tracker) GetUpdateChan()chan bool{
	t.updateChan = make(chan bool)
	return t.updateChan

}

func (t *tracker) Update(instKey string, childKey string, childVal int64){
	t.mux.Lock()

		select{
		case val := <- t.writeChan:
			log.Info("WRITE CHAN////////////////////////////",val)
			log.Info("Storing the offset in the following instances:", t.Cached.Instances)
			t.storage.Write(t.Cached)

		case key := <-t.removeChan:
			log.Info("REMOVE CHAN////////////////////////////", key)
			//delete(t.Cached.Instances, key)
			log.Info("[Remove]Current status : ", t.Cached.Instances)

		case res := <-t.updateChan:
			log.Info("UPDATE CHAN////////////////////////////", res)
			if _, instance_exist := t.Cached.Instances[instKey]; instance_exist{
				tailNode := t.Cached.Instances[instKey]
				if _, childNode_exist := tailNode[childKey]; childNode_exist{

					tailNode[childKey] = childVal
					t.Cached.Instances[instKey] = tailNode
					//log.Info("[Update]Current status : ", t.Cached.Instances)
				}

			}
		default:
			// do some other crap here
		}

	t.mux.Unlock()
	runtime.Gosched()
	
	//log.Info("---------------------------------------------------------")
	//t.LoadTailers()
	//log.Info("---------------------------------------------------------")
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
