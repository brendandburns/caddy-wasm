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
* C - https://github.com/brendandburns/caddy-wasm/tree/main/examples/c
* Dotnet - https://github.com/dev-wasm/dev-wasm-dotnet/blob/main/www-wasi/Program.cs

More to come in the future!

## Version loading
The wasm loader supports loading multiple versions of the same functionality and then chosing between them using a header.

To see this in action run:
```sh
# Run caddy pointed at the releases at https://github.com/brendandburns/caddy-wasm/releases
$ xcaddy run --config examples/github-versions.json
# Get the default version
$ curl localhost:2015
This is the beta release!
# Select a version
$ curl -H "x-caddy-wasm-version: 0.0.1" localhost:2015
This is the alpha release!
```

The specific configuration for this looks like:
```json
{
    "handle": [{
        "handler": "wasm",
        "wasm_file": "tinygo.wasm",
        "wasm_github": "brendandburns/caddy-wasm"
    }]
}
```
The `wasm_file` directive specifies the file to look for in the release assets, releases without that file will be skipped.
The `wasm_github` directive points to the repository where the releases are contained.

### Other version loaders
Currently there is also a glob loader that supports loading all files that match a glob:
```json
{
    "handle": [{
        "handler": "wasm",
        "wasm_file": "test/*.wasm"
    }]
}
```

All files matching the glob (from the working directory where `caddy` is launched) will be served.
Their "version" is the complete file path that matches the glob.

## Bugs/Features
I'm sure there are many, please feel free to file issues.