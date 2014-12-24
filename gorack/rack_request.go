package gorack

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

type RackRequest struct {
	Request     *http.Request
	SERVER_NAME string
	SERVER_PORT string
}

func NewRackRequest(r *http.Request, serverName, serverPort string) *RackRequest {
	return &RackRequest{
		Request:     r,
		SERVER_NAME: serverName,
		SERVER_PORT: serverPort,
	}
}

func (rr *RackRequest) Bytes() []byte {

	r := rr.Request

	type kval struct{ k, val string }

	items := []kval{
		{"REQUEST_METHOD", r.Method},
		{"SCRIPT_NAME", r.URL.Path},
		{"PATH_INFO", r.URL.Path},
		{"QUERY_STRING", r.URL.RawQuery},
		{"SERVER_NAME", rr.SERVER_NAME},
		{"SERVER_PORT", rr.SERVER_PORT},
	}

	for k, vals := range r.Header {
		items = append(items, kval{k, strings.Join(vals, "; ")})
	}

	buf := &bytes.Buffer{}

	for _, item := range items {
		buf.WriteString(fmt.Sprintf("%s: %s%s", item.k, item.val, delim))
	}

	buf.WriteString(delim)

	return buf.Bytes()
}
