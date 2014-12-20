package gorack

import "net/http"

type RackRequest struct {
	REQUEST_METHOD string
	SCRIPT_NAME    string
	PATH_INFO      string
	QUERY_STRING   string
	SERVER_NAME    string
	SERVER_PORT    string
	HTTP_vars      []string
}

func NewRackRequest(r *http.Request, serverName, serverPort string) RackRequest {
	return RackRequest{
		REQUEST_METHOD: r.Method,
		SCRIPT_NAME:    r.URL.Path,
		PATH_INFO:      r.URL.Path,
		QUERY_STRING:   r.URL.RawQuery,
		SERVER_NAME:    serverName,
		SERVER_PORT:    serverPort,
	}
}
