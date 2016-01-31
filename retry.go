package greq

import (
	"github.com/cenkalti/backoff"
	"net/http"
	"time"
)

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
