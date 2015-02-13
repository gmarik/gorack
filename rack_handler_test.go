package gorack

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"runtime"
	"strings"
	"sync"
	"testing"
)

var testBody = "OMG\x00, test body\nhalps"

// var echoBody = `{"REQUEST_METHOD"=>"POST", "SCRIPT_NAME"=>"", "PATH_INFO"=>"/", "QUERY_STRING"=>"", "SERVER_NAME"=>"host", "SERVER_PORT"=>"port", "HTTP_ACCEPT_ENCODING"=>"gzip", "HTTP_CONTENT_LENGTH"=>"21", "HTTP_CONTENT_TYPE"=>"text/plain", "HTTP_USER_AGENT"=>"Go 1.1 package http"}` + testBody

var echoBody = `{"HTTP_ACCEPT_ENCODING"=>"gzip", "HTTP_CONTENT_LENGTH"=>"21", "HTTP_CONTENT_TYPE"=>"text/plain", "HTTP_USER_AGENT"=>"Go 1.1 package http", "PATH_INFO"=>"/", "QUERY_STRING"=>"", "REQUEST_METHOD"=>"POST", "SCRIPT_NAME"=>"", "SERVER_NAME"=>"host", "SERVER_PORT"=>"port"}` + testBody

func TestRackHandler(t *testing.T) {

	var cases = []struct{ in, exp, script string }{
		{testBody, testBody, "./ruby/test/config_test.ru"},
		{testBody, echoBody, "./ruby/test/echo.ru"},
	}

	for _, v := range cases {
		exp := v.exp
		got, host, port, err := submit(v.in, v.script)

		if err != nil {
			t.Error(err)
		}

		if echoBody == v.exp {
			exp = strings.Replace(v.exp, `=>"port"`, `=>"`+port+`"`, -1)
			exp = strings.Replace(exp, `=>"host"`, `=>"`+host+`"`, -1)
		}

		if exp != got {
			t.Errorf("\nExp:%s\nGot:%s", exp, got)
		}
	}
}

func submit(body string, rackScript string) (string, string, string, error) {

	handler := NewRackHandler(rackScript)

	if err := handler.StartRackProcess(); err != nil {
		return "", "", "", err
	}

	ts := httptest.NewServer(handler)

	log.Println(ts.URL)

	defer handler.StopRackProcess()
	defer ts.Close()

	res, err := http.Post(ts.URL, "text/plain", strings.NewReader(body))
	if err != nil {
		return "", "", "", err
	}

	got, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {

		return "", "", "", err
	}

	u, _ := url.Parse(ts.URL)

	host, port, _ := net.SplitHostPort(u.Host)

	return string(got), host, port, nil
}

func BenchmarkRackHandler(b *testing.B) {

	// runtime.GOMAXPROCS(2)

	handler := NewRackHandler("./ruby/test/echo.ru")

	if err := handler.StartRackProcess(); err != nil {
		b.Error(err)
	}

	ts := httptest.NewServer(handler)

	log.Println(ts.URL)

	defer handler.StopRackProcess()
	defer ts.Close()

	req := make(chan int)

	wg := sync.WaitGroup{}

	worker := func(i int) {
		for rn := range req {
			wg.Add(1)

			body := fmt.Sprintf("worker %d:%d", i, rn)
			// log.Println(body)
			_, err := http.Post(ts.URL, "text/plain", strings.NewReader(body))
			if err != nil {
				b.Error(err)
			}
			wg.Done()
		}
	}

	for i := 0; i < 10; i += 1 {
		go worker(i)
	}

	b.ResetTimer()

	printed := true
	for i := 0; i < 100; {
		select {
		case req <- i:
			i += 1
			printed = false
		default:
			if !printed {
				log.Println("Waiting")
				printed = true
			}
			runtime.Gosched()
		}
	}
	wg.Wait()
}
