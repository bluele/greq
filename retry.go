package greq

import (
	"github.com/cenkalti/backoff"
	"net/http"
	"time"
)

// Retry retry to request at even intervals.
// retry: retry number
// interval: retry interval
func Retry(retry int, interval time.Duration) RequestHandler {
	return func(req *Request, doReq RequestMethod) (*http.Response, error) {
		var res *http.Response
		_retry := retry
		return res, retryInterval(func() error {
			var err error
			res, err = doReq(req)
			if err != nil && _retry > 0 {
				_retry--
				return err
			}
			return nil
		}, interval)
	}
}

func retryInterval(cb func() error, interval time.Duration) error {
	for {
		if err := cb(); err == nil {
			return err
		}
		time.Sleep(interval)
	}
	return nil
}

// Exponential backoff
// retry: retry number
// b: cenkalti backoff object
func RetryBackoff(retry int, b backoff.BackOff) RequestHandler {
	return func(req *Request, doReq RequestMethod) (*http.Response, error) {
		var res *http.Response
		_retry := retry
		return res, backoff.Retry(func() error {
			var err error
			res, err = doReq(req)
			if err != nil && _retry > 0 {
				_retry--
				return err
			}
			return nil
		}, b)
	}
}

func NewBackOff() backoff.BackOff {
	return backoff.NewExponentialBackOff()
}

// We should retry if specified function returns true.
// cb: callback function after request. If this function returns true, retry request cancelled.
// interval: retry number
func RetryOnResult(cb func(*http.Response, error) bool, interval time.Duration) RequestHandler {
	return func(req *Request, doReq RequestMethod) (*http.Response, error) {
		for {
			res, err := doReq(req)
			if cb(res, err) {
				return res, err
			}
			if interval > 0 {
				time.Sleep(interval)
			}
		}
	}
}
