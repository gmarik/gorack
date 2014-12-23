package gorack

import (
	"bytes"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

// delimiter used to separate headers from one another
// and body from the rest of the headers
const delim = "\x00"

type RackResponse struct {
	rackResponse io.Reader
	Headers      map[string][]string
	StatusCode   int
	Body         io.Reader
}

func NewResponse(r io.Reader) *RackResponse {
	return &RackResponse{rackResponse: r}
}

func (r *RackResponse) Parse() error {
	headerBuffer := &bytes.Buffer{}

	// while determining headers size
	// read(tee) headers into separate
	// buffer for futher processing
	reader := io.TeeReader(r.rackResponse, headerBuffer)

	// at some point reader reader reaches the body
	r.Body = r.rackResponse

	// read char by char to correctly land at body start
	buf := make([]byte, 1, 1)

	// end of headers, end of line
	eol, eoh := false, false

	for {
		_, err := reader.Read(buf)

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		char := string(buf[0]) // single char

		// eoh marks end of headers
		eoh = eol && delim == char

		if eoh {
			break
		}

		eol = delim == char
	}

	if err := r.parseHeaders(headerBuffer); err != nil {
		return err
	}

	return nil
}

func (r *RackResponse) parseHeaders(buf io.Reader) error {
	headers, err := ioutil.ReadAll(buf)

	if err != nil {
		return err
	}

	lines := bytes.Split(headers, []byte(delim))

	// first header is a status code
	code, err := strconv.Atoi(string(lines[0]))

	if err != nil {
		return err
	}

	r.StatusCode = code

	r.Headers = make(map[string][]string)

	for _, line := range lines[1:] {
		hdr := string(line)

		if hdr == "" {
			continue
		}

		kvs := strings.SplitN(hdr, ": ", 2)
		r.Headers[kvs[0]] = strings.Split(kvs[1], "; ")
	}

	return nil
}
