package gorack

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
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

hello world!`

	response2 = `200
X-This: a messsage
Content-Length: 5

hello`
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

	buf := &bytes.Buffer{}

	// copy body into buffer
	// just uses "streaming"
	io.Copy(buf, r.Body)

	body, err := ioutil.ReadAll(buf)

	if err != nil {
		t.Error(err)
	}

	exp := "hello world!"

	if !reflect.DeepEqual(string(body), exp) {
		t.Errorf("\nExp %s\nGot %s", exp, string(body))
	}
}

func TestReponse2(t *testing.T) {
	read, write, err := os.Pipe()

	write.Write([]byte(response2))
	write.Close()

	r := NewResponse(read)

	if err := r.Parse(); err != nil {
		t.Error(err)
	}

	_, err = ioutil.ReadAll(r.Body)

	if err != nil {
		t.Error(err)
	}
}
