package gorack

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestRackRequest(t *testing.T) {

	url, err := url.Parse("http://addre.ss/path/script.ext?query=param1")

	if err != nil {
		t.Error(err)
	}

	r := &http.Request{Method: "GET", URL: url}

	got := NewRackRequest(r, "server", "port")

	exp := &RackRequest{
		REQUEST_METHOD: "GET",
		SCRIPT_NAME:    "/path/script.ext",
		PATH_INFO:      "/path/script.ext",
		QUERY_STRING:   "query=param1",
		SERVER_NAME:    "server",
		SERVER_PORT:    "port",
	}

	if !reflect.DeepEqual(got, exp) {
		t.Errorf("\nExp: %v\nGot: %v", exp, got)
	}
}
