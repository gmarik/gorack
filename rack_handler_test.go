package gorack

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var testBody = "OMG\x00, test body\nhalps"
var echoBody = `{"REQUEST_METHOD"=>"POST", "SCRIPT_NAME"=>"/", "PATH_INFO"=>"/", "QUERY_STRING"=>"", "SERVER_NAME"=>"server", "SERVER_PORT"=>"port", "Accept-Encoding"=>"gzip", "Content-Length"=>"21", "Content-Type"=>"text/plain", "User-Agent"=>"Go 1.1 package http"}` + testBody

func TestRackHandler(t *testing.T) {

	var cases = []struct{ in, exp, script string }{
		{testBody, testBody, "./ruby/test/config_test.ru"},
		{testBody, echoBody, "./ruby/test/echo.ru"},
	}

	for _, v := range cases {
		exp, got := v.exp, submit(v.in, v.script, t)

		if !bytes.Equal([]byte(exp), []byte(got)) {
			t.Errorf("\nGot:%s\nExp:%s", []byte(got), []byte(exp))
		}
	}
}

func submit(body string, rackScript string, t *testing.T) string {
	// package variable
	GorackRunner = "./ruby/libexec/gorack"

	ts := httptest.NewServer(NewRackHandler(rackScript))
	defer ts.Close()

	res, err := http.Post(ts.URL, "text/plain", strings.NewReader(body))
	if err != nil {
		t.Error(err)
	}

	got, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		t.Error(err)
	}

	return string(got)
}
