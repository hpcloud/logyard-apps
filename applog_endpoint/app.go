package applog_endpoint

import (
	"stackato/client"
	"stackato/server"
)

func recentLogs(token, appGUID string, num int) ([]string, error) {
	endpoint := server.GetClusterConfig().Endpoint
	targetUrl := "https://" + endpoint
	space := "" // we don't care about space
	cli := client.NewRestClient(targetUrl, token, space)
	return cli.GetLogsRaw(appGUID, num)
}
