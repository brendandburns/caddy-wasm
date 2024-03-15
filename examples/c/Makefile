wasi_sdk := /usr/local/lib/wasi-sdk-21.0
cc := ${wasi_sdk}/bin/clang

.phony: clean

default: server.wasm

libwasihttp.a:
	wget https://github.com/dev-wasm/dev-wasm-c/releases/download/0.2.0/libwasihttp-0.2.0.tgz
	tar -xvzf libwasihttp-0.2.0.tgz
	rm libwasihttp-0.2.0.tgz
	cat wasi_http.h | grep -v '#include "client.h"' > wasi_http.h.fixed
	mv wasi_http.h.fixed wasi_http.h

%.o : %.c
	@echo "Compiling c file into o file"
	${cc} -c $< -o $@

server.wasm: libwasihttp.a server.o; ${cc} server.o -L. -lwasihttp -o server.wasm

clean: ; rm -f *.o *.wasm *.a *.tgz wasi_http.h
