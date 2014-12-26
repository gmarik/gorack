package gorack

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"
)

var (
	testUrl = "http://addre.ss/path/script.ext?query=param1"

	testRequest = &RackRequest{
		Request:     nil,
		SERVER_NAME: "server",
		SERVER_PORT: "port",
	}

	testRequestString = "" +
		"REQUEST_METHOD: GET" + delim +
		"SCRIPT_NAME: /path/script.ext" + delim +
		"PATH_INFO: /path/script.ext" + delim +
		"QUERY_STRING: query=param1" + delim +
		"SERVER_NAME: server" + delim +
		"SERVER_PORT: port" + delim +
		"Accept-Encoding: gzip, deflate" + delim +
		"Accept-Language: da, en-gb; q=0.8, en" + delim +
		"Connection: keep-alive" + delim +
		delim

	headers = map[string][]string{
		"Accept-Encoding": {"gzip, deflate"},
		"Accept-Language": {"da, en-gb", "q=0.8, en"},
		"Connection":      {"keep-alive"},
	}
)

func TestRackRequestBytesSerialization(t *testing.T) {

	url, err := url.Parse(testUrl)

	if err != nil {
		t.Error(err)
	}

	r := &http.Request{Method: "GET", URL: url, Header: headers}

	testRequest.Request = r
	got := testRequest.Bytes()

	exp := []byte(testRequestString)

	if !bytes.Equal(got, exp) {
		t.Errorf("\nExp: %s\nGot: %s", exp, got)
	}
}
