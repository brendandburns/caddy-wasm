package caddy_wasm

import (
	"fmt"
	"net/http"
	"os"

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
	wasi   *wasi_http.WasiHTTP
	rt     wazero.Runtime
	loader *WebAssemblyLoader

	WebAssemblyFile string `json:"wasm_file,omitempty"`
}

func (WebAssembly) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.wasm",
		New: func() caddy.Module { return new(WebAssembly) },
	}
}

func (w *WebAssembly) Validate() error {
	if len(w.WebAssemblyFile) == 0 {
		return fmt.Errorf("no wasm file specified")
	}
	if _, err := os.Stat(w.WebAssemblyFile); os.IsNotExist(err) {
		return fmt.Errorf("wasm file does not exist: %s", w.WebAssemblyFile)
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

	var data []byte
	data, err = os.ReadFile(w.WebAssemblyFile)
	if err != nil {
		return err
	}

	var mod wazero.CompiledModule
	mod, err = w.rt.CompileModule(ctx, data)
	if err != nil {
		return err
	}

	w.loader = &WebAssemblyLoader{
		rt:     w.rt,
		module: mod,

		// TODO: fill out this config here
		moduleConfig: wazero.NewModuleConfig(),
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
	instance, err := w.loader.GetOrLoad(req.Context())
	defer instance.Release()
	if err != nil {
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
