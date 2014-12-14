package gorack

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strconv"
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

const MAX_HEADERS_SIZE = 512 * 16 * 1024 // 256 * 16K

type RackResponse struct {
	rackResponse io.Reader
	Headers      map[string][]string
	Status       int
	Body         io.Reader

	buf         *bytes.Buffer
	headersSize uint
}

func (r *RackResponse) Parse() error {
	r.buf = &bytes.Buffer{}

	// while determining headers size
	// read(tee) headers into separate
	// buffer for futher processing
	reader := bufio.NewReader(io.TeeReader(r.rackResponse, r.buf))

	// end of headers, end of line
	eol, eoh := false, false

	// at some point reader reader reaches the body
	r.Body = reader

	// read char by char to correctly land at body start
	char := make([]byte, 1, 1)

	for {
		n, err := reader.Read(char)

		if err != nil {
			return err
		}

		r.headersSize += uint(n)

		// \n\n marks end of headers
		eoh = eol && byte('\n') == char[0]

		if eoh {
			break
		}

		eol = byte('\n') == char[0]
	}

	fmt.Println("Read ", r.headersSize, " bytes")

	err := r.ParseHeaders()

	if err != nil {
		return err
	}

	return nil
}

func (r *RackResponse) ParseHeaders() error {
	var delim = byte('\n')

	// reads headers based on previously determined r.headersSize
	headers := make([]byte, r.headersSize, r.headersSize)
	_, err := r.buf.Read(headers)

	if err != nil {
		return err
	}

	fmt.Println(string(headers))

	if err != nil {
		return err
	}

	lines := bytes.Split(headers, []byte{delim})

	// first header is a status code
	code, err := strconv.Atoi(string(lines[0]))

	if err != nil {
		return err
	}

	r.Status = code

	r.Headers = make(map[string][]string)

	for _, line := range lines[1:] {
		hdr := string(line)
		// fmt.Println(hdr)

		if hdr == "" {
			continue
		}

		kvs := strings.SplitN(hdr, ": ", 2)
		r.Headers[kvs[0]] = strings.Split(kvs[1], "; ")
	}

	return nil
}

func NewReader(r io.Reader) *RackResponse {
	return &RackResponse{rackResponse: r}
}

func TestResponseParse(t *testing.T) {
	r := NewReader(strings.NewReader(response))

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

	if r.Body == nil {
		t.Errorf("r.Body can't be nil")
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		t.Error(err)
	}

	if reflect.DeepEqual(string(body), "hello world!") {
		t.Errorf("\nGot %s", string(body))
	}
}
