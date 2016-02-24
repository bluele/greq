package greq_test

import (
	"encoding/json"
	"fmt"
	"github.com/bluele/greq"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Response struct {
	ID   int
	Name string
}

func TestGetRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	}))
	defer ts.Close()

	var err error

	_, err = greq.Get(ts.URL).Do()
	if err != nil {
		t.Error(err)
		return
	}

	if v, err := greq.Get(ts.URL).String(); err != nil {
		t.Error(err)
	} else if v != "ok" {
		t.Error(`response should be "ok"`)
	}

	if v, err := greq.Get(ts.URL).Bytes(); err != nil {
		t.Error(err)
	} else if string(v) != "ok" {
		t.Error(`response should be "ok"`)
	}
}

func TestPostRequest(t *testing.T) {
	var (
		expectedKey   = "key"
		expectedValue = "value"
	)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			fmt.Fprint(w, "failed")
			return
		}
		r.ParseForm()
		fmt.Fprintf(w, r.Form.Get(expectedKey))
	}))
	defer ts.Close()

	res, err := greq.Post(ts.URL, []byte(expectedKey+"="+expectedValue)).Do()
	if err != nil {
		t.Error(err)
		return
	}
	if res == nil {
		t.Error("res should not be nil")
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
		return
	}
	if string(body) != expectedValue {
		t.Errorf("body should not be %v", string(body))
	}
}

func TestJSONResponse(t *testing.T) {
	correctResponse := Response{1, "bluele"}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, _ := json.Marshal(correctResponse)
		fmt.Fprint(w, string(resp))
	}))
	defer ts.Close()

	resp := &Response{}
	err := greq.Get(ts.URL).JSON(resp)
	if err != nil {
		t.Error(err)
		return
	}

	if resp.ID != correctResponse.ID || resp.Name != correctResponse.Name {
		t.Errorf("response should be %v", correctResponse)
	}
}
