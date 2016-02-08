package main

import (
	"errors"
	"fmt"
	"github.com/bluele/greq"
	"io/ioutil"
	"net/http"
)

func main() {
	res, err := greq.Get("http://example.com/notfound.html").
		RequestHandler(greq.RetryBackoff(3, greq.NewBackOff())).
		ResponseHandler(func(res *http.Response, err error) error {
		if res != nil && res.StatusCode >= 400 && res.StatusCode < 500 {
			return errors.New("40X error")
		}
		return err
	}).Do()
	if err != nil {
		fmt.Println("error:", err.Error())
		return
	}
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
}
