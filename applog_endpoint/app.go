package main

import (
	"encoding/json"
	"fmt"
	"stackato/client"
	"stackato/server"
)

func recentLogs(token, appGUID string) string {
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
