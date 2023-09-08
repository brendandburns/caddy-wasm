# Example Caddy WASM module

## Prerequisites
This requires `go` and `tinygo` installed.

## Building
```sh
tinygo build -o tinygo.wasm -target wasi
```

## Running
```
# Go back to the base caddy-wasm directory
CADDY_WASM_DIR=/some/path/to/caddy-wasm
cd ${CADDY_WASM_DIR}
xcaddy run --config examples/tinygo/caddy.json
```