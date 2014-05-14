package main

import (
	"fmt"
	"github.com/ActiveState/log"
	"github.com/ActiveState/logyard-apps/applog_endpoint/wsutil"
	"github.com/gorilla/mux"
	"net/http"
)

func sendRecent(stream *wsutil.WebSocketStream, args *Arguments) error {
	if args.Num <= 0 {
		// First authorize with the CC by fetching something
		_, err := recentLogs(args.Token, args.GUID, 1)
		if err != nil {
			return err
		}
	} else {
		// Recent history requested?
		recentLogs, err := recentLogs(args.Token, args.GUID, args.Num)
		if err != nil {
			return err
		}
		for _, line := range recentLogs {
			err = stream.Send(line)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func recentHandlerWs(
	w http.ResponseWriter, r *http.Request, stream *wsutil.WebSocketStream) {
	args, err := ParseArguments(r)
	if err != nil {
		stream.Fatalf("Invalid arguments; %v", err)
		return
	}

	if err := sendRecent(stream, args); err != nil {
		stream.Fatalf("%v", err)
		return
	}
}

func tailHandlerWs(
	w http.ResponseWriter, r *http.Request, stream *wsutil.WebSocketStream) {
	args, err := ParseArguments(r)
	if err != nil {
		stream.Fatalf("Invalid arguments; %v", err)
		return
	}

	if err := sendRecent(stream, args); err != nil {
		stream.Fatalf("%v", err)
		return
	}

	drain, err := NewAppLogDrain(args.GUID)
	if err != nil {
		stream.Fatalf("Unable to create drain: %v", err)
		return
	}
	ch, err := drain.Start()
	if err != nil {
		stream.Fatalf("Unable to start drain: %v", err)
	}

	err = stream.Forward(ch)
	if err != nil {
		log.Infof("%v", err)
		drain.Stop(err)
	}

	// We expect drain.Wait to not block at this point.
	if err := drain.Wait(); err != nil {
		if _, ok := err.(wsutil.WebSocketStreamError); !ok {
			log.Warnf("Error from app log drain server: %v", err)
		}
	}
}

func serve() error {
	addr := fmt.Sprintf(":%d", PORT)
	r := mux.NewRouter()
	r.Handle("/v2/apps/{guid}/tail",
		wsutil.WebSocketHandler(wsutil.HandlerFunc(tailHandlerWs)))
	r.Handle("/v2/apps/{guid}/recent",
		wsutil.WebSocketHandler(wsutil.HandlerFunc(recentHandlerWs)))

	http.Handle("/", r)
	return http.ListenAndServe(addr, nil)
}
