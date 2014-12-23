package gorack

import (
	"bytes"
	"fmt"
	"net/http"
)

type RackRequest struct {
	REQUEST_METHOD string
	SCRIPT_NAME    string
	PATH_INFO      string
	QUERY_STRING   string
	SERVER_NAME    string
	SERVER_PORT    string
	HTTP_vars      []string
}

func NewRackRequest(r *http.Request, serverName, serverPort string) *RackRequest {
	return &RackRequest{
		REQUEST_METHOD: r.Method,
		SCRIPT_NAME:    r.URL.Path,
		PATH_INFO:      r.URL.Path,
		QUERY_STRING:   r.URL.RawQuery,
		SERVER_NAME:    serverName,
		SERVER_PORT:    serverPort,
	}
}

func (r *RackRequest) Bytes() []byte {
	items := []struct {
		k, val string
	}{
		{"REQUEST_METHOD", r.REQUEST_METHOD},
		{"SCRIPT_NAME", r.SCRIPT_NAME},
		{"PATH_INFO", r.PATH_INFO},
		{"QUERY_STRING", r.QUERY_STRING},
		{"SERVER_NAME", r.SERVER_NAME},
		{"SERVER_PORT", r.SERVER_PORT},
	}

	buf := &bytes.Buffer{}

	for _, item := range items {
		buf.WriteString(fmt.Sprintf("%s: %s%s", item.k, item.val, delim))
	}

	buf.WriteString(delim)

	return buf.Bytes()
}
