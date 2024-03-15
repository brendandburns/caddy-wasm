# Example Caddy WASM module

## Prerequisites
The wasi-sdk must be installed (`clang`)

## Building
```sh
make server.wasm
```

## Running
```
# Go back to the base caddy-wasm directory
CADDY_WASM_DIR=/some/path/to/caddy-wasm
cd ${CADDY_WASM_DIR}
xcaddy run --config examples/c/caddy.json
```