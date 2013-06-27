package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"net/url"
	"io/ioutil"
)

func TestCallWalkCompiler2(t *testing.T) {
	Tempdir = "tmp"
	PathToWalkCompiler2 = "WALK_Compiler2.exe"
	
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Compile(w, r)
	}))
	defer ts.Close()
	
	values := make(url.Values)
	values.Set("doc", "{\"Type\":\"posxml\", \"Code\": \"<waitkey/>\", \"Output\": \"\"}")
	r, _ := http.PostForm(ts.URL, values)
	body, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	
	returnCompiler:= "{\"err\": null,\"posxml\":{ \"base64\": \"dQ0=\", \"size\": 2, \"integers\": 0, \"strings\": 0, \"maxvars\": 512, \"functions\": 0, \"maxfuncs\": 128}}\r\n"

	if string(body) != returnCompiler {
		t.Errorf("got %s want json with code compiled", body)
	}
}