package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"net/url"
	"io/ioutil"
)

func TestCallEcho(t *testing.T) {
	Tempdir = "tmp"
	
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Compile(w, r)
	}))
	defer ts.Close()
	
	values := make(url.Values)
	values.Set("doc", "{\"Type\":\"echo\", \"Code\": \"abc\", \"Output\": \"\"}")
	r, _ := http.PostForm(ts.URL, values)
	body, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	if string(body) != "abc" {
		t.Errorf("got %s want abc", body)
	}
}