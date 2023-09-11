package caddy_wasm

import (
	"net/http"
	"net/http/httptest"
	"reflect"
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

func TestGithub(t *testing.T) {
	v, e := ForGithubRepository("brendandburns", "caddy-wasm", "tinygo.wasm")
	if e != nil {
		t.Errorf("unexpected error: %v", e.Error())
	}
	vs, e := v.GetVersions()
	if e != nil {
		t.Errorf("unexpected error: %v", e.Error())
	}
	if !reflect.DeepEqual(vs, []string{"0.0.1"}) {
		t.Errorf("unexpected output: %v", vs)
	}

	data, e := v.GetWebAssembly("0.0.1")
	if e != nil {
		t.Errorf("unexpected error: %v", e.Error())
	}
	if len(data) == 0 {
		t.Errorf("unexpected data: %v", data)
	}
}
