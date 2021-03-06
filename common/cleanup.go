package common

import (
	"github.com/ActiveState/log"
	"github.com/ActiveState/tail"
	"os"
	"os/signal"
	"syscall"
)

func cleanup() {
	log.Info("cleanup: closing open inotify watches")
	tail.Cleanup()
}

// Fatal is like log.Fatal, but invokes cleanup (tail) before exiting.
func Fatal(format string, v ...interface{}) {
	log.Fatal0(format, v...)
	cleanup()
	os.Exit(1)
}

func RegisterTailCleanup() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	for sig := range c {
		log.Warnf("captured signal %v; exiting after cleanup", sig)
		cleanup()
		os.Exit(1)
	}
}
