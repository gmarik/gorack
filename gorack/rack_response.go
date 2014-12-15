package gorack

import (
	"bytes"
	"io"
	"strconv"
	"strings"
)

type RackResponse struct {
	rackResponse io.Reader
	Headers      map[string][]string
	StatusCode   int
	Body         io.Reader

	buf         *bytes.Buffer
	headersSize uint
}

func NewResponse(r io.Reader) *RackResponse {
	return &RackResponse{rackResponse: r}
}

func (r *RackResponse) Parse() error {
	r.buf = &bytes.Buffer{}

	// while determining headers size
	// read(tee) headers into separate
	// buffer for futher processing
	reader := io.TeeReader(r.rackResponse, r.buf)

	// at some point reader reader reaches the body
	r.Body = r.rackResponse

	// read char by char to correctly land at body start
	char := make([]byte, 1, 1)

	// end of headers, end of line
	eol, eoh := false, false

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

	// fmt.Println("Read ", r.headersSize, " bytes")

	if err := r.parseHeaders(); err != nil {
		return err
	}

	return nil
}

func (r *RackResponse) parseHeaders() error {
	var delim = byte('\n')

	// reads headers based on previously determined r.headersSize
	headers := make([]byte, r.headersSize, r.headersSize)
	_, err := r.buf.Read(headers)

	if err != nil {
		return err
	}

	lines := bytes.Split(headers, []byte{delim})

	// first header is a status code
	code, err := strconv.Atoi(string(lines[0]))

	if err != nil {
		return err
	}

	r.StatusCode = code

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
