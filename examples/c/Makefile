wasi_sdk := /usr/local/lib/wasi-sdk-20.0
cc := ${wasi_sdk}/bin/clang
ar := ${wasi_sdk}/bin/ar
lib_objs := client.o wasi_http.o

.phony: gen clean run

default: server.wasm

# This is necessary for a full `make` from clean b/c client.c doesn't exist when the make starts...
# there must be a better way to do this.
client.o: gen
	${cc} -c client.c -o client.o

%.o : %.c
	@echo "Compiling c file into o file"
	${cc} -c $< -o $@

libwasihttp.a: ${lib_objs}
	${ar} cr libwasihttp.a ${lib_objs}

wasi-http: ; git clone https://github.com/WebAssembly/wasi-http; cd wasi-http; git checkout v0.2.0-rc-2023-11-10 ; cp ../client.wit wit/client.wit; cd ../

gen: wasi-http ; wit-bindgen c ./wasi-http/wit -w client

server.wasm: gen server.o libwasihttp.a; ${cc} server.o -L. -lwasihttp -o server.wasm

clean: ; rm -f client.c client.h *.o *.wasm *.a; rm -rf wasi-http
