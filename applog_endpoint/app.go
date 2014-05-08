package main

import (
	"stackato/client"
	"stackato/server"
)

func recentLogs(token, appGUID string, num int) ([]client.AppLogLine, error) {
	endpoint := server.GetClusterConfig().Endpoint
	targetUrl := "https://" + endpoint
	space := "" // we don't care about space
	cli := client.NewRestClient(targetUrl, token, space)
	return cli.GetLogs(appGUID, num)
}
