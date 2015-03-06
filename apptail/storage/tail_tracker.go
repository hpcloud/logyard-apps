package storage

import (
	"runtime"
	"fmt"
	"github.com/ActiveState/log"
	"os"
	"sync"
)

type Tracker interface {
	Update(key string, tailNode *TailNode)
	LoadTailers() *Tailer
	Remove(key string)
}

type TailNode map[string]string

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

func (t *tracker) LoadTailers() *Tailer {
	t.mux.Lock()
	t.storage.Load(&t.Cached)
	log.Info("loaded the following tailers from previous session:", t.Cached.Instances)
	t.mux.Unlock()
	runtime.Gosched()
	return t.Cached
}

func (t *tracker) Update(key string, tailNode *TailNode) {
	log.Info("Inserting the following key:", key)
	t.mux.Lock()
	t.Cached.IsLive = true
	t.Cached.Instances[key] = (*tailNode)
	t.storage.Write(t.Cached)
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
