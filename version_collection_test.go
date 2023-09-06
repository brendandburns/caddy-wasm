package caddy_wasm

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
)

type testHandler struct {
	requests []*http.Request
}

func (t *testHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	t.requests = append(t.requests, req)
	res.WriteHeader(200)
	res.Write([]byte("Ok"))
}

func TestUrlVersions(t *testing.T) {
	handler := &testHandler{}
	server := httptest.NewServer(handler)

	collection := ForURL(server.URL)

	v, e := collection.GetVersions()
	if e != nil {
		t.Errorf("Unexpected error: %v", e.Error())
	}
	if len(v) != 1 {
		t.Errorf("Unexpected versions: %v", v)
	}
	d, e := collection.GetWebAssembly(v[0])
	if e != nil {
		t.Errorf("Unexpected error: %v", e.Error())
	}
	if string(d) != "Ok" {
		t.Errorf("Unexpected data: %v", string(d))
	}
}

func TestFileVersions(t *testing.T) {
	m := fstest.MapFS{}
	m["some/dir/foo.wasm"] = &fstest.MapFile{Data: []byte("foo")}
	m["some/dir/bar.wasm"] = &fstest.MapFile{Data: []byte("bar")}
	f := &fileSystemVersionCollection{
		m,
		"some/dir/",
		"*.wasm",
	}
	v, e := f.GetVersions()
	if e != nil {
		t.Errorf("Unexpected error: %v", e.Error())
	}
	tests := []struct {
		v string
		d string
	}{
		{"foo.wasm", "foo"},
		{"bar.wasm", "bar"},
	}
	if len(v) != len(tests) {
		t.Errorf("Unexpected versions: %v", v)
	}
	for _, test := range tests {
		d, e := f.GetWebAssembly(test.v)
		if e != nil {
			t.Errorf("Unexpected error: %v", e.Error())
		}
		if string(d) != test.d {
			t.Errorf("Unexpected data: %v", string(d))
		}
	}
}
