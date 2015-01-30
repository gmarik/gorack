package gorack

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sort"
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
		{"SCRIPT_NAME", ""}, // TODO: buil properly
		{"PATH_INFO", r.URL.Path},
		{"QUERY_STRING", r.URL.RawQuery},
		{"SERVER_NAME", rr.SERVER_NAME},
		{"SERVER_PORT", rr.SERVER_PORT},
	}

	// sort keys so order is predictable
	keys := make([]string, 0)
	for k, _ := range r.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		items = append(items, kval{http_header(k), strings.Join(r.Header[k], "; ")})
	}

	buf := &bytes.Buffer{}

	for _, item := range items {
		buf.WriteString(fmt.Sprintf("%s: %s%s", item.k, item.val, delim))
	}

	buf.WriteString(delim)

	return buf.Bytes()
}

func (r *RackRequest) WriteTo(w io.WriteCloser) error {
	if _, err := w.Write(r.Bytes()); err != nil {
		return err
	}

	if _, err := io.Copy(w, r.Request.Body); err != nil {
		return err
	}

	w.Close()

	return nil
}

func http_header(name string) string {
	name = strings.Replace(name, "-", "_", -1)
	return "HTTP_" + strings.ToUpper(name)
}
