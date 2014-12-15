package gorack

import (
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

const (
	response = `200
Server: nginx/1.6.0
Content-Type: text/html
Content-Length: 0
Last-Modified: Sat, 06 Sep 2014 15:37:58 GMT
Date: Wed, 26 Nov 2014 23:49:32 GMT
Connection: keep-alive
Set-Cookie: UserID=JohnDoe; Max-Age=3600; Version=1

hello world!
`
)

func TestResponseParse(t *testing.T) {
	r := NewResponse(strings.NewReader(response))

	result := map[string][]string{
		"Server":         []string{"nginx/1.6.0"},
		"Content-Type":   []string{"text/html"},
		"Content-Length": []string{"0"},
		"Last-Modified":  []string{"Sat, 06 Sep 2014 15:37:58 GMT"},
		"Date":           []string{"Wed, 26 Nov 2014 23:49:32 GMT"},
		"Connection":     []string{"keep-alive"},
		"Set-Cookie":     []string{"UserID=JohnDoe", "Max-Age=3600", "Version=1"},
	}

	if err := r.Parse(); err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(r.Headers, result) {
		t.Errorf("\nExp %s,\nGot %s", result, r.Headers)
	}

	if r.StatusCode != 200 {
		t.Error("invalid response code")
	}

	if r.Body == nil {
		t.Errorf("response Body can't be nil")
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		t.Error(err)
	}

	if reflect.DeepEqual(string(body), "hello world!") {
		t.Errorf("\nGot %s", string(body))
	}
}
