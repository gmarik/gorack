package gorack

import (
	"bufio"
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

type RackResponse struct {
	rackResponse io.Reader
	Headers      map[string][]string
	Status       int
	Body         io.Reader
}

func (r *RackResponse) Parse() error {
	scanner := bufio.NewScanner(r.rackResponse)
	scanner.Split(bufio.ScanLines)
	scanner.Scan()
	status := scanner.Text()

	code, err := strconv.Atoi(status)

	if err != nil {
		return nil
	}

	r.Status = code

	r.Headers = make(map[string][]string)

	// r.bodyReaderOffset = scanner.Pos.Offset
	r.Body = r.rackResponse

	for scanner.Scan() {
		entry := scanner.Text()

		fmt.Println(entry)
		// reached body?
		if entry == "" {
			return nil
		}

		// k: val; val
		line := strings.SplitN(entry, ": ", 2)
		vals := strings.Split(line[1], "; ")

		r.Headers[line[0]] = vals
	}

	// if err := scanner.Err(); err != nil {
	//             fmt.Fprintln(os.Stderr, "There was an error with the scanner in attached container", err)
	//         }

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

	r.Parse()

	if !reflect.DeepEqual(r.Headers, result) {
		t.Errorf("\nExp %s,\nGot %s", result, r.Headers)
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		t.Error(err)
	}

	if string(body) != "hello world" {
		t.Errorf("\nGot %s", string(body))
	}
}
