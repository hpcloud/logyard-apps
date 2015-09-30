package applog_endpoint

import (
	"testing"

	"github.com/hpcloud/logyard-apps/applog_endpoint/config"
)

func TestGetApplogEndpointUriDefault(t *testing.T) {
	config.GetClusterConfig().Endpoint = "api.stackato.example"
	config.GetConfig().Hostname = ""
	uri := getApplogEndpointUri()
	expected := "logs.stackato.example"
	if uri != expected {
		t.Errorf("Got unexpected applog endpoint uri %v, expected %v",
			uri, expected)
	}
}

func TestGetApplogEndpointUriCustom(t *testing.T) {
	config.GetClusterConfig().Endpoint = "api.stackato.example"
	expected := "example.test"
	config.GetConfig().Hostname = expected
	uri := getApplogEndpointUri()
	if uri != expected {
		t.Errorf("Got unexpected applog endpoint uri %v, expected %v",
			uri, expected)
	}
}
