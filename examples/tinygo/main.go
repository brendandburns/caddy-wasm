package main

import (
	"net/http"

	"github.com/dev-wasm/dev-wasm-go/http/server"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("This is the beta release!"))
	})
	server.ListenAndServe(nil)
}
