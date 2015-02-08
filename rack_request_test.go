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
		"HTTP_ACCEPT_ENCODING: gzip, deflate" + delim +
		"HTTP_ACCEPT_LANGUAGE: da, en-gb; q=0.8, en" + delim +
		"HTTP_CONNECTION: keep-alive" + delim +
		"PATH_INFO: /path/script.ext" + delim +
		"QUERY_STRING: query=param1" + delim +
		"REQUEST_METHOD: GET" + delim +
		"SCRIPT_NAME: " + delim +
		"SERVER_NAME: server" + delim +
		"SERVER_PORT: port" + delim +
		delim

	testHeaders = map[string][]string{
		"Accept-Encoding": {"gzip, deflate"},
		"Accept-Language": {"da, en-gb", "q=0.8, en"},
		"Connection":      {"keep-alive"},
	}
)

func TestRackRequestHeaderSerialization(t *testing.T) {

	url, err := url.Parse(testUrl)

	if err != nil {
		t.Error(err)
	}

	r := &http.Request{Method: "GET", URL: url, Header: testHeaders}

	testRequest.Request = r

	buf := &bytes.Buffer{}

	testRequest.writeHeaders(buf, testRequest.headers())

	got := buf.Bytes()

	exp := []byte(testRequestString)

	if !bytes.Equal(got, exp) {
		t.Errorf("\nExp: %s\nGot: %s", exp, got)
	}
}
