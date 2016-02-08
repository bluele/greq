package greq

import (
	"net/http"
)

// Default request handler
func defaultRequestHandler(req *Request, doReq RequestMethod) (*http.Response, error) {
	return doReq(req)
}
