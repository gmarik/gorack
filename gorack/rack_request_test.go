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

	testRequestString = `REQUEST_METHOD: GET
SCRIPT_NAME: /path/script.ext
PATH_INFO: /path/script.ext
QUERY_STRING: query=param1
SERVER_NAME: server
SERVER_PORT: port`
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
	exp := []byte(testRequestString + "\n\n")

	got := testRequest.Bytes()

	if !bytes.Equal(got, exp) {
		t.Errorf("\nExp: %s\nGot: %s", exp, got)
	}

}
