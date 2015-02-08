package gorack

import (
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

func (rr *RackRequest) headers() http.Header {
	r := rr.Request

	headers := http.Header{
		"SCRIPT_NAME":    []string{""}, // TODO: build properly
		"REQUEST_METHOD": []string{r.Method},
		"PATH_INFO":      []string{r.URL.Path},
		"QUERY_STRING":   []string{r.URL.RawQuery},
		"SERVER_NAME":    []string{rr.SERVER_NAME},
		"SERVER_PORT":    []string{rr.SERVER_PORT},
	}

	for k, val := range r.Header {
		headers[http_header(k)] = val
	}

	return headers
}

func (rr *RackRequest) writeHeaders(out io.Writer, headers http.Header) error {

	// sort keys so order is predictable
	keys := make([]string, 0)

	for k, _ := range headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, sort_key := range keys {
		k, val := sort_key, headers[sort_key]
		fmt.Fprintf(out, "%s: %s%s", k, strings.Join(val, "; "), delim)
	}

	_, err := out.Write([]byte(delim))

	return err
}

func (r *RackRequest) WriteTo(w io.WriteCloser) error {

	if err := r.writeHeaders(w, r.headers()); err != nil {
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
