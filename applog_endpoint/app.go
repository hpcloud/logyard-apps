package main

import (
	"stackato/client"
	"stackato/server"
)

func recentLogs(token, appGUID string) ([]client.AppLogLine, error) {
	endpoint := server.GetClusterConfig().Endpoint
	targetUrl := "https://" + endpoint
	cli := client.NewRestClient(targetUrl, token, "")
	return cli.GetLogs(appGUID, 5)
}
