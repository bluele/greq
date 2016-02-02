package greq

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func Bytes(res *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func String(res *http.Response, err error) (string, error) {
	body, err := Bytes(res, err)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func JSON(res *http.Response, err error, ptr interface{}) error {
	body, err := Bytes(res, err)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, ptr)
}
