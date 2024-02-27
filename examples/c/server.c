#include "wasi_http.h"

#include <stdio.h>

void handler_fn(wasi_http_request_t *req, wasi_http_response_t *res) {
    res->status_code = 200;
    sprintf(res->body, "{ \"msg\": \"Hello world!\", \"path\": \"%s\", \"authority\": \"%s\"}\n", req->path_query, req->authority);
}

int main() {
    handler = handler_fn;
}

// This is hack and doesn't actually do anything for this server 
bool exports_wasi_cli_run_run() {
    return !main();
}
