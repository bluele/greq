package greq

import (
	"github.com/cenkalti/backoff"
	"net/http"
	"time"
)

// Retry retry to request at even intervals.
// retry: retry number
// interval: retry interval
func Retry(retry int, interval time.Duration) func(*Request, func() (*http.Response, error)) error {
	return func(req *Request, doReq func() (*http.Response, error)) error {
		_retry := retry
		return retryInterval(func() error {
			_, err := doReq()
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
func RetryBackoff(retry int, b backoff.BackOff) func(*Request, func() (*http.Response, error)) error {
	return func(req *Request, doReq func() (*http.Response, error)) error {
		_retry := retry
		return backoff.Retry(func() error {
			_, err := doReq()
			if err != nil && _retry > 0 {
				_retry--
				return err
			}
			return nil
		}, b)
	}
}

// We should retry, specified function returns true.
// cb: callback function after request. If this function returns true, retry request cancelled.
// interval: retry number
func RetryOnResult(cb func(*http.Response, error) bool, interval time.Duration) func(*Request, func() (*http.Response, error)) error {
	return func(req *Request, doReq func() (*http.Response, error)) error {
		for {
			res, err := doReq()
			if cb(res, err) {
				return err
			}
			if interval > 0 {
				time.Sleep(interval)
			}
		}
	}
}
