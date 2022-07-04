package router_toy

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testdata = []struct {
	sourcePath string
	targetPath string
	handle     http.HandlerFunc
	hopeVal    string
}{
	{
		"/aa/:name/c",
		"/aa/qwe/c",
		func(writer http.ResponseWriter, request *http.Request) {
			writer.Write([]byte("aaa"))
		},
		"aaa",
	},
	{
		"/b",
		"/b",
		func(writer http.ResponseWriter, request *http.Request) {
			writer.Write([]byte("bbb"))
		},
		"bbb",
	},
	{
		"/c/*filepath",
		"/c/qew.txt",
		func(writer http.ResponseWriter, request *http.Request) {
			writer.Write([]byte("ccc"))
		},
		"ccc",
	},
}

func OpenServer() *Router {
	r := New()
	for _, v := range testdata {
		r.Get(v.sourcePath, v.handle)
	}
	return r
}

func TestRouter(t *testing.T) {
	r := OpenServer()
	for _, val := range testdata {
		req := httptest.NewRequest(http.MethodGet, val.targetPath, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		res := w.Result()
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}
		if string(data) != val.hopeVal {
			t.Errorf("expected %v got %v", val.hopeVal, string(data))
		}
		res.Body.Close()
	}
}
