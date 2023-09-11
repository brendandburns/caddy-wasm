package caddy_wasm

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/stealthrocket/wasi-go/imports/wasi_http"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// Derived from the starting example here: https://caddyserver.com/docs/extending-caddy

func init() {
	caddy.RegisterModule(WebAssembly{})
	httpcaddyfile.RegisterHandlerDirective("visitor_ip", parseCaddyfile)
}

type WebAssembly struct {
	wasi     *wasi_http.WasiHTTP
	rt       wazero.Runtime
	versions VersionCollection
	loader   *VersionedLoader

	WebAssemblyFile   string `json:"wasm_file,omitempty"`
	WebAssemblyURL    string `json:"wasm_url,omitempty"`
	WebAssemblyGithub string `json:"wasm_github,omitempty"`
}

func (WebAssembly) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.wasm",
		New: func() caddy.Module { return new(WebAssembly) },
	}
}

func (w *WebAssembly) Validate() error {
	if len(w.WebAssemblyFile) == 0 && len(w.WebAssemblyURL) == 0 && len(w.WebAssemblyGithub) == 0 {
		return fmt.Errorf("no wasm file or url specified")
	}
	if len(w.WebAssemblyFile) > 0 && len(w.WebAssemblyURL) > 0 {
		return fmt.Errorf("both wasm file and url specified")
	}

	if len(w.WebAssemblyFile) > 0 && len(w.WebAssemblyGithub) == 0 {
		files, err := filepath.Glob(w.WebAssemblyFile)
		if err != nil {
			return err
		}
		if len(files) == 0 {
			return fmt.Errorf("wasm file does not exist: %s", w.WebAssemblyFile)
		}
	}
	if len(w.WebAssemblyURL) > 0 {
		_, err := url.Parse(w.WebAssemblyURL)
		if err != nil {
			return fmt.Errorf("invalid wasm url: %s", w.WebAssemblyURL)
		}
	}
	if len(w.WebAssemblyGithub) > 0 {
		parts := strings.Split(w.WebAssemblyGithub, "/")
		if len(parts) != 2 {
			return fmt.Errorf("expected <owner>/<repo> for github")
		}
		if len(w.WebAssemblyFile) == 0 {
			return fmt.Errorf("wasm_file is required for github")
		}
	}
	return nil
}

func (w *WebAssembly) GetVersionCollection() VersionCollection {
	if len(w.WebAssemblyFile) > 0 && len(w.WebAssemblyGithub) == 0 {
		dir := filepath.Dir(w.WebAssemblyFile)
		glob := filepath.Base(w.WebAssemblyFile)
		return ForFilesystemGlob(dir, glob)
	}
	if len(w.WebAssemblyURL) > 0 {
		return ForURL(w.WebAssemblyURL)
	}
	if len(w.WebAssemblyGithub) > 0 {
		parts := strings.Split(w.WebAssemblyGithub, "/")
		fmt.Printf("%s %s %s\n", parts[0], parts[1], w.WebAssemblyFile)
		v, _ := ForGithubRepository(parts[0], parts[1], w.WebAssemblyFile)
		return v
	}
	return nil
}

func (w *WebAssembly) Provision(ctx caddy.Context) error {
	config := wazero.NewRuntimeConfig().
		WithCloseOnContextDone(true)
	var err error
	err = nil
	w.rt = wazero.NewRuntimeWithConfig(ctx, config)
	defer func() {
		if err != nil {
			w.rt.Close(ctx)
		}
	}()

	w.versions = w.GetVersionCollection()
	if w.versions == nil {
		return fmt.Errorf("no wasm file or url specified")
	}

	w.loader, err = MakeVersionedLoader(ctx, w.versions, w.rt)
	if err != nil {
		return err
	}

	_, err = wasi_snapshot_preview1.Instantiate(ctx, w.rt)
	if err != nil {
		return err
	}

	w.wasi = wasi_http.MakeWasiHTTP()
	if err = w.wasi.Instantiate(ctx, w.rt); err != nil {
		return err
	}

	return nil
}

func (w *WebAssembly) ServeHTTP(res http.ResponseWriter, req *http.Request, next caddyhttp.Handler) error {
	version := req.Header.Get("x-caddy-wasm-version")
	if len(version) == 0 {
		version = w.loader.DefaultVersion
	}
	instance, err := w.loader.GetOrLoad(req.Context(), version)
	if instance == nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Failed to load wasm: " + err.Error()))
		return fmt.Errorf("no wasm instance available")
	}
	defer instance.Release()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("Failed to load wasm: " + err.Error()))
		return err
	}
	handler := w.wasi.MakeHandler(instance.module)
	handler.ServeHTTP(res, req)
	return nil
}

func (w *WebAssembly) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if !d.Args(&w.WebAssemblyFile) {
			return d.ArgErr()
		}
	}
	return nil
}

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var w WebAssembly
	err := w.UnmarshalCaddyfile(h.Dispenser)
	return &w, err
}

// Interface guards
var (
	_ caddy.Provisioner           = (*WebAssembly)(nil)
	_ caddy.Validator             = (*WebAssembly)(nil)
	_ caddyhttp.MiddlewareHandler = (*WebAssembly)(nil)
	_ caddyfile.Unmarshaler       = (*WebAssembly)(nil)
)
