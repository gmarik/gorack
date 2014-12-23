package gorack

import (
	"bytes"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

var (
	testUrl = "http://addre.ss/path/script.ext?query=param1"

	testRequest = &RackRequest{
		REQUEST_METHOD: "GET",
		SCRIPT_NAME:    "/path/script.ext",
		PATH_INFO:      "/path/script.ext",
		QUERY_STRING:   "query=param1",
		SERVER_NAME:    "server",
		SERVER_PORT:    "port",
	}

	testRequestString = "" +
		"REQUEST_METHOD: GET\x00" +
		"SCRIPT_NAME: /path/script.ext\x00" +
		"PATH_INFO: /path/script.ext\x00" +
		"QUERY_STRING: query=param1\x00" +
		"SERVER_NAME: server\x00" +
		"SERVER_PORT: port\x00"
)

func TestRackRequest(t *testing.T) {
	url, err := url.Parse(testUrl)

	if err != nil {
		t.Error(err)
	}

	r := &http.Request{Method: "GET", URL: url}

	exp := testRequest
	got := NewRackRequest(r, "server", "port")

	if !reflect.DeepEqual(got, exp) {
		t.Errorf("\nExp: %v\nGot: %v", exp, got)
	}
}

func TestRackRequestBytesSerialization(t *testing.T) {
	exp := []byte(testRequestString + delim)

	got := testRequest.Bytes()

	if !bytes.Equal(got, exp) {
		t.Errorf("\nExp: %s\nGot: %s", exp, got)
	}
}

func TestRequestHeaders(t *testing.T) {
	t.SkipNow()
}
