package docker_events

import (
	"encoding/json"
	"github.com/ActiveState/log"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Event struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	From   string `json:"from"`
	Time   int64  `json:"time"`
}

func Stream() chan *Event {
	ch := make(chan *Event)
	res := getDockerEvents(3)

	go func() {
		defer res.Body.Close()

		// Read the streaming json from the events endpoint
		// http://docs.docker.io/en/latest/api/docker_remote_api_v1.3/#monitor-docker-s-events
		d := json.NewDecoder(res.Body)
		for {
			var event Event
			if err := d.Decode(&event); err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal(err)
			}
			ch <- &event
		}

		close(ch)
	}()

	return ch
}

func getDockerEvents(retries int) *http.Response {
	c := http.Client{}
	for attempt := 0; attempt < retries; attempt++ {
		res, err := c.Get("http://localhost:4243/events")
		if err != nil {
			if (attempt + 1) == retries {
				log.Fatalf("Failed to read from docker daemon; giving up retrying: %v", err)
			}
			log.Warnf("Docker connection error (%v); retrying after 1 second.", err)
			time.Sleep(time.Second)
		} else {
			return res
		}
	}
	panic("unreachable")
}

type Docker struct {
	Id string `json:"id"`
}

func GetLiveDockerContainers(retries int) map[string]bool {
	const ID_LENGTH = 12
	allDockerIds := make(map[string]bool)
	var httpResByte []byte
	c := http.Client{}
	for attempt := 0; attempt < retries; attempt++ {
		res, err := c.Get("http://localhost:4243/containers/json")
		defer res.Body.Close()
		if err != nil {
			if (attempt + 1) == retries {
				log.Fatalf("Failed to read from docker daemon; giving up retrying: %v", err)
			}
			log.Warnf("Docker connection error (%v); retrying after 1 second.", err)
			time.Sleep(time.Second)
		} else {
			httpResByte, err = ioutil.ReadAll(res.Body)
			if err != nil {
				panic(err)

			}
			var jsonData []Docker
			err = json.Unmarshal([]byte(httpResByte), &jsonData)

			if err != nil {
				panic(err)

			}

			for _, element := range jsonData {
				shortenedKey := element.Id[:ID_LENGTH]
				allDockerIds[shortenedKey] = true
			}

		}
	}
	return allDockerIds
}
