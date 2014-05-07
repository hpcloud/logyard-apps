package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"stackato/client"
	"stackato/server"
)

func serveDemo(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	result := demo(token, r.FormValue("appGUID"))
	w.Write([]byte(result))
}

func demo(token, appGUID string) string {
	endpoint := server.GetClusterConfig().Endpoint
	targetUrl := "https://" + endpoint
	cli := client.NewRestClient(targetUrl, token, "")
	logs, err := cli.GetLogs(appGUID, 5)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	data, err := json.Marshal(logs)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	return string(data)
}
