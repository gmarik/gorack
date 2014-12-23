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

	// while determining headers size
	// read(tee) headers into separate
	// buffer for futher processing
	headerBuffer := &bytes.Buffer{}
	reader := io.TeeReader(r.rackResponse, headerBuffer)

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

	code, headers, err := parseHeaders(headerBuffer)

	if err != nil {
		return err
	}

	r.StatusCode = code
	r.Headers = headers
	r.Body = r.rackResponse

	return nil
}

func parseHeaders(buf io.Reader) (int, map[string][]string, error) {
	headers, err := ioutil.ReadAll(buf)

	if err != nil {
		return 0, nil, err
	}

	lines := bytes.Split(headers, []byte(delim))

	// first header is a status code
	code, err := strconv.Atoi(string(lines[0]))

	if err != nil {
		return 0, nil, err
	}

	hdrs := make(map[string][]string)

	for _, line := range lines[1:] {
		hdr := string(line)

		if hdr == "" {
			continue
		}

		kvs := strings.SplitN(hdr, ": ", 2)
		hdrs[kvs[0]] = strings.Split(kvs[1], "; ")
	}

	return code, hdrs, nil
}
