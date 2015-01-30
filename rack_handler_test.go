package gorack

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

var testBody = "OMG\x00, test body\nhalps"
var echoBody = `{"REQUEST_METHOD"=>"POST", "SCRIPT_NAME"=>"", "PATH_INFO"=>"/", "QUERY_STRING"=>"", "SERVER_NAME"=>"host", "SERVER_PORT"=>"port", "Accept-Encoding"=>"gzip", "Content-Length"=>"21", "Content-Type"=>"text/plain", "User-Agent"=>"Go 1.1 package http"}` + testBody

func TestRackHandler(t *testing.T) {

	var cases = []struct{ in, exp, script string }{
		{testBody, testBody, "./ruby/test/config_test.ru"},
		{testBody, echoBody, "./ruby/test/echo.ru"},
	}

	for _, v := range cases {
		exp := v.exp
		got, host, port := submit(v.in, v.script, t)

		if echoBody == v.exp {
			exp = strings.Replace(v.exp, `=>"port"`, `=>"`+port+`"`, -1)
			exp = strings.Replace(exp, `=>"host"`, `=>"`+host+`"`, -1)
		}

		if exp != got {
			t.Errorf("\nGot:%s\nExp:%s", got, exp)
		}
	}
}

func submit(body string, rackScript string, t *testing.T) (string, string, string) {
	// package variable
	GorackRunner = "./ruby/libexec/gorack"

	handler := NewRackHandler(rackScript)

	if err := handler.StartRackProcess(); err != nil {
		t.Error(err)
	}

	ts := httptest.NewServer(handler)

	log.Println(ts.URL)

	defer handler.StopRackProcess()
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

	u, _ := url.Parse(ts.URL)

	host, port, _ := net.SplitHostPort(u.Host)

	return string(got), host, port
}
