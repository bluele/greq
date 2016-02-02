package greq_test

import (
	"encoding/json"
	"fmt"
	"github.com/bluele/greq"
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
	}

	if resp.ID != correctResponse.ID || resp.Name != correctResponse.Name {
		t.Errorf("response should be %v", correctResponse)
	}
}
