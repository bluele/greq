package greq

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httputil"
	"sync"
)

type Request struct {
	method  string
	rawurl  string
	header  http.Header
	body    []byte
	client  *http.Client
	cookies []*http.Cookie
	once    sync.Once

	responseHandler ResponseHandler
	requestHandler  RequestHandler

	debug bool
}

type (
	RequestHandler  func(func() (*http.Response, error)) error
	ResponseHandler func(*http.Response, error) error
)

func Get(rawurl string) *Request {
	return New("GET", rawurl)
}

func Post(rawurl string) *Request {
	return New("POST", rawurl)
}

func Put(rawurl string) *Request {
	return New("PUT", rawurl)
}

func Delete(rawurl string) *Request {
	return New("DELETE", rawurl)
}

func New(method, rawurl string) *Request {
	req := &Request{}
	req.method = method
	req.rawurl = rawurl
	return req
}

func defaultRequestHandler(doReq func() (*http.Response, error)) error {
	_, err := doReq()
	return err
}

func (req *Request) Client() *http.Client {
	return req.client
}

func (req *Request) SetClient(client http.Client) *Request {
	req.client = &client
	return req
}

func (req *Request) Header() http.Header {
	return req.header
}

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

func (req *Request) AddHeader(key string, values ...string) *Request {
	for _, value := range values {
		req.header.Add(key, value)
	}
	return req
}

func (req *Request) SetBody(body []byte) *Request {
	req.body = body
	return req
}

func (req *Request) SetUseragent(value string) *Request {
	req.SetHeader("User-Agent", value)
	return req
}

func (req *Request) Do() (*http.Response, error) {
	req.once.Do(func() {
		if req.client == nil {
			req.client = &(*http.DefaultClient)
		}
	})
	var (
		res *http.Response
		err error
	)
	rh := req.requestHandler
	if rh == nil {
		rh = defaultRequestHandler
	}
	if e := rh(func() (*http.Response, error) {
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

func (req *Request) RequestHandler(requestHandler func(*Request, func() (*http.Response, error)) error) *Request {
	req.requestHandler = func(doReq func() (*http.Response, error)) error {
		if err := requestHandler(req, doReq); err != nil {
			return err
		}
		return nil
	}
	return req
}

func (req *Request) ResponseHandler(handler func(res *http.Response, err error) error) *Request {
	req.responseHandler = handler
	return req
}

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

func (req *Request) String() (string, error) {
	return String(req.Do())
}

func (req *Request) Bytes() ([]byte, error) {
	return Bytes(req.Do())
}

func (req *Request) JSON(ptr interface{}) error {
	res, err := req.Do()
	return JSON(res, err, ptr)
}

func (req *Request) Debug(debug bool) *Request {
	req.debug = debug
	return req
}
