# WebAssembly extension for Caddy
This is a very experimental work in progress, use at your own risk!

PRs welcome!

## Running
```sh
xcaddy run --config caddy.json
curl localhost:2015
```

## Building WASM modules
The server uses the [WASI-HTTP](https://github.com/WebAssembly/wasi-http) specification.

There are examples here:
* Go (TinyGo) - https://github.com/brendandburns/caddy-wasm/tree/main/examples/tinygo
* Dotnet - https://github.com/dev-wasm/dev-wasm-dotnet/blob/main/www-wasi/Program.cs
* C - https://github.com/stealthrocket/wasi-go/blob/main/testdata/c/http/server.c

More to come in the future!

## Bugs/Features
I'm sure there are many, please feel free to file issues.