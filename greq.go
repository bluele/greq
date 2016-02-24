package greq

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httputil"
)

// Request context object.
type Request struct {
	method  string
	rawurl  string
	header  http.Header
	body    []byte
	client  *http.Client
	cookies []*http.Cookie

	responseHandler ResponseHandler
	requestHandler  RequestHandler

	debug bool
}

type (
	RequestMethod   func(*Request) (*http.Response, error)
	RequestHandler  func(*Request, RequestMethod) (*http.Response, error)
	ResponseHandler func(*http.Response, error) error
)

const (
	contentType = "Content-type"

	defaultPOSTContentType = "application/x-www-form-urlencoded"
)

var (
	Debug = false
)

// Set get method
func Get(rawurl string) *Request {
	return New().Get().SetURL(rawurl)
}

// Set post method
func Post(rawurl string, body []byte) *Request {
	return New().Post().SetURL(rawurl).SetBody(body).SetHeader(contentType, defaultPOSTContentType)
}

// Create a new request object.
func New() *Request {
	req := &Request{}
	req.header = make(http.Header)
	req.debug = Debug
	req.client = http.DefaultClient
	return req
}

// Method returns method name.
func (req *Request) Method() string {
	return req.method
}

// SetMethod sets specified string as http method.
func (req *Request) SetMethod(method string) *Request {
	req.method = method
	return req
}

// Set get method
func (req *Request) Get() *Request {
	return req.SetMethod("GET")
}

// Set post method
func (req *Request) Post() *Request {
	return req.SetMethod("POST")
}

// URL returns url string.
func (req *Request) URL() string {
	return req.rawurl
}

// SetURL sets specified string as request url.
func (req *Request) SetURL(rawurl string) *Request {
	req.rawurl = rawurl
	return req
}

// Client returns current *http.Client
func (req *Request) Client() *http.Client {
	return req.client
}

// SetClient sets *http.Client
func (req *Request) SetClient(client http.Client) *Request {
	req.client = &client
	return req
}

// Header returns current http.Header.
func (req *Request) Header() http.Header {
	return req.header
}

// SetHeader sets key-values as request header.
func (req *Request) SetHeader(key string, values ...string) *Request {
	for i, value := range values {
		if i == 0 {
			req.header.Set(key, value)
		} else {
			req.header.Add(key, value)
		}
	}
	return req
}

// SetHeader adds key-values to request header.
func (req *Request) AddHeader(key string, values ...string) *Request {
	for _, value := range values {
		req.header.Add(key, value)
	}
	return req
}

// SetBody sets specified body as request body.
func (req *Request) SetBody(body []byte) *Request {
	req.body = body
	return req
}

// SetUseragent sets a specified string as request useragent.
func (req *Request) SetUseragent(value string) *Request {
	req.SetHeader("User-Agent", value)
	return req
}

// Do HTTP requests using itself parameters.
func (req *Request) Do() (*http.Response, error) {
	var (
		res *http.Response
		err error
	)
	rh := req.requestHandler
	if rh == nil {
		rh = defaultRequestHandler
	}
	if res, e := rh(req, func(req *Request) (*http.Response, error) {
		res, err = req.doReq(req.method, req.rawurl)
		if req.responseHandler != nil {
			err = req.responseHandler(res, err)
		}
		return res, err
	}); e != nil {
		return res, e
	}
	if err != nil {
		return res, err
	}
	return res, nil
}

// RequestHandler hooks an event which before sending request.
func (req *Request) RequestHandler(handler RequestHandler) *Request {
	req.requestHandler = handler
	return req
}

// ResponseHandler hooks an event which after sending request.
func (req *Request) ResponseHandler(handler ResponseHandler) *Request {
	req.responseHandler = handler
	return req
}

// AddCookie adds a cookie to request headers.
func (req *Request) AddCookie(cookie *http.Cookie) *Request {
	req.cookies = append(req.cookies, cookie)
	return req
}

func (req *Request) doReq(method, rawurl string) (*http.Response, error) {
	r, err := http.NewRequest(method, rawurl, bytes.NewBuffer(req.body))
	if err != nil {
		return nil, err
	}
	r.Header = req.header
	for _, c := range req.cookies {
		r.AddCookie(c)
	}

	if req.debug {
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Println(err)
		}
		log.Println(string(dump))
	}

	res, err := req.client.Do(r)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Type converter for string response body.
func (req *Request) String() (string, error) {
	return String(req.Do())
}

// Type converter for []byte response body.
func (req *Request) Bytes() ([]byte, error) {
	return Bytes(req.Do())
}

// JSON bind a response body to specified object.
func (req *Request) JSON(ptr interface{}) error {
	res, err := req.Do()
	return JSON(res, err, ptr)
}

// Give true argument, print debug log when do request.
func (req *Request) Debug(debug bool) *Request {
	req.debug = debug
	return req
}
