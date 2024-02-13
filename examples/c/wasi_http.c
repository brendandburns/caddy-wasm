#include <assert.h>
#include <string.h>
#include <stdio.h>
#include "client.h"
#include "wasi_http.h"

int main();

handler_func handler = NULL;

void wasi_http_set_handler(handler_func handler_fn) {
    handler = handler_fn;
}

char* c_str(client_string_t* str) {
    char* ptr = malloc(str->len + 1);
    sprintf(ptr, "%s", str->ptr);
    return ptr;
}

void exports_wasi_http_0_2_0_rc_2023_11_10_incoming_handler_handle(exports_wasi_http_0_2_0_rc_2023_11_10_incoming_handler_own_incoming_request_t request, exports_wasi_http_0_2_0_rc_2023_11_10_incoming_handler_own_response_outparam_t response_out) {
    // This is a hack, that is required for now for initialization.
    main();
    client_string_t path_query;
    client_string_t authority;

    wasi_http_0_2_0_rc_2023_11_10_types_method_incoming_request_path_with_query(wasi_http_0_2_0_rc_2023_11_10_types_borrow_incoming_request(request), &path_query);
    char* path_query_cstr = c_str(&path_query);
    wasi_http_0_2_0_rc_2023_11_10_types_method_incoming_request_authority(wasi_http_0_2_0_rc_2023_11_10_types_borrow_incoming_request(request), &authority);
    char* authority_cstr = c_str(&authority);


    char buff[64 * 1024];
    wasi_http_request_t wasi_request = {
        .path_query = path_query_cstr,
        .authority = authority_cstr,
    };
    wasi_http_response_t wasi_response = {
        .body = buff,
        .body_max_len = 64 * 1024,
    };
    if (handler != NULL) {
        handler(&wasi_request, &wasi_response);
    }

    wasi_http_0_2_0_rc_2023_11_10_types_tuple2_field_key_field_value_t headers [] = {
        {
            .f0 = {
                .ptr = (uint8_t*) "Content-type",
                .len = strlen("Content-type"),
            },
            .f1 = {
                .ptr = (uint8_t*) "text/plain",
                .len = strlen("text/plain"),
            },
        },
    };

    wasi_http_0_2_0_rc_2023_11_10_types_list_tuple2_field_key_field_value_t header_list = {
        .ptr = headers,
        .len = 1,
    };

    wasi_http_0_2_0_rc_2023_11_10_types_own_headers_t f;
    wasi_http_0_2_0_rc_2023_11_10_types_header_error_t err;

    wasi_http_0_2_0_rc_2023_11_10_types_static_fields_from_list(&header_list, &f, &err);

    wasi_http_0_2_0_rc_2023_11_10_types_own_outgoing_response_t res = wasi_http_0_2_0_rc_2023_11_10_types_constructor_outgoing_response(f);
    wasi_http_0_2_0_rc_2023_11_10_types_method_outgoing_response_set_status_code(wasi_http_0_2_0_rc_2023_11_10_types_borrow_outgoing_response(res), wasi_response.status_code);

    wasi_http_0_2_0_rc_2023_11_10_types_own_outgoing_body_t body;
    wasi_http_0_2_0_rc_2023_11_10_types_method_outgoing_response_body(wasi_http_0_2_0_rc_2023_11_10_types_borrow_outgoing_response(res), &body);

    wasi_http_0_2_0_rc_2023_11_10_types_result_own_outgoing_response_error_code_t res_err;
    res_err.is_err = false;
    res_err.val.ok = res;
    wasi_http_0_2_0_rc_2023_11_10_types_static_response_outparam_set(response_out, &res_err);

    wasi_http_0_2_0_rc_2023_11_10_types_own_output_stream_t stream;
    wasi_http_0_2_0_rc_2023_11_10_types_method_outgoing_body_write(wasi_http_0_2_0_rc_2023_11_10_types_borrow_outgoing_body(body), &stream);

    wasi_io_0_2_0_rc_2023_11_10_streams_list_u8_t contents;
    contents.ptr = (uint8_t*) wasi_response.body;
    contents.len = strlen((const char*) contents.ptr);
    wasi_io_0_2_0_rc_2023_11_10_streams_stream_error_t stream_err;
    wasi_io_0_2_0_rc_2023_11_10_streams_method_output_stream_blocking_write_and_flush(wasi_io_0_2_0_rc_2023_11_10_streams_borrow_output_stream(stream), &contents, &stream_err);

    wasi_io_0_2_0_rc_2023_11_10_streams_output_stream_drop_own(stream);

    wasi_http_0_2_0_rc_2023_11_10_types_error_code_t error_code;
    wasi_http_0_2_0_rc_2023_11_10_types_static_outgoing_body_finish(body, NULL, &error_code);

    free(path_query_cstr);
    free(authority_cstr);
}

void free_response(wasi_http_response_t* response) {
    for (int i = 0; i < response->headers.len; i++) {
        //free(response->headers.headers[i].name);
        free(response->headers.headers[i].value);
    }
    free(response->headers.headers);
}

int wasi_http_request(uint8_t method_tag, uint8_t scheme_tag, const char * authority_str, const char* path_query_str, const char* body, wasi_http_response_t* response_out) {
    wasi_http_0_2_0_rc_2023_11_10_types_tuple2_field_key_field_value_t content_type[] = {{
        .f0 = { .ptr = (uint8_t*)"User-agent", .len = 10 },
        .f1 = { .ptr = (uint8_t*)"WASI-HTTP/0.0.1", .len = 15},
    },
    {
        .f0 = { .ptr = (uint8_t*)"Content-type", .len = 12 },
        .f1 = { .ptr = (uint8_t*)"application/json", .len = 16},
    }};
    wasi_http_0_2_0_rc_2023_11_10_types_list_tuple2_field_key_field_value_t headers_list = {
        .ptr = &content_type[0],
        .len = 2,
    };
    wasi_http_0_2_0_rc_2023_11_10_types_own_fields_t headers;
    wasi_http_0_2_0_rc_2023_11_10_types_header_error_t err;
    if (!wasi_http_0_2_0_rc_2023_11_10_types_static_fields_from_list(&headers_list, &headers, &err)) {
        printf("Header create failed\n");
        return 8;
    }
    wasi_http_0_2_0_rc_2023_11_10_types_method_t method = { .tag = method_tag };
    wasi_http_0_2_0_rc_2023_11_10_types_scheme_t scheme = { .tag = scheme_tag };
    client_string_t path_query, authority;
    client_string_set(&path_query, (char*) path_query_str);
    client_string_set(&authority, (char *) authority_str);

    wasi_http_0_2_0_rc_2023_11_10_types_own_outgoing_body_t out_body;

    wasi_http_0_2_0_rc_2023_11_10_types_own_outgoing_request_t req = wasi_http_0_2_0_rc_2023_11_10_types_constructor_outgoing_request(headers);
    bool ok = wasi_http_0_2_0_rc_2023_11_10_types_method_outgoing_request_set_method(wasi_http_0_2_0_rc_2023_11_10_types_borrow_outgoing_request(req), &method) &&
        wasi_http_0_2_0_rc_2023_11_10_types_method_outgoing_request_set_path_with_query(wasi_http_0_2_0_rc_2023_11_10_types_borrow_outgoing_request(req), &path_query) &&
        wasi_http_0_2_0_rc_2023_11_10_types_method_outgoing_request_set_scheme(wasi_http_0_2_0_rc_2023_11_10_types_borrow_outgoing_request(req), &scheme) &&
        wasi_http_0_2_0_rc_2023_11_10_types_method_outgoing_request_set_authority(wasi_http_0_2_0_rc_2023_11_10_types_borrow_outgoing_request(req), &authority) &&
        wasi_http_0_2_0_rc_2023_11_10_types_method_outgoing_request_body(wasi_http_0_2_0_rc_2023_11_10_types_borrow_outgoing_request(req), &out_body);

    if (!ok) {
        printf("Error creating request\n");
        return 4;
    }
    if (body != NULL) {
        wasi_http_0_2_0_rc_2023_11_10_types_own_output_stream_t ret;
        if (!wasi_http_0_2_0_rc_2023_11_10_types_method_outgoing_body_write(wasi_http_0_2_0_rc_2023_11_10_types_borrow_outgoing_body(out_body), &ret)) {
            printf("Error getting output stream\n");
            return 7;
        }
        wasi_io_0_2_0_rc_2023_11_10_streams_list_u8_t buf = {
            .ptr = (uint8_t *) body,
            .len = strlen(body),
        };
        wasi_io_0_2_0_rc_2023_11_10_streams_stream_error_t stream_err;
        // TODO check error here.
        wasi_io_0_2_0_rc_2023_11_10_streams_method_output_stream_blocking_write_and_flush(wasi_io_0_2_0_rc_2023_11_10_streams_borrow_output_stream(ret), &buf, &stream_err);
        wasi_io_0_2_0_rc_2023_11_10_streams_output_stream_drop_own(ret);
    }

    wasi_http_0_2_0_rc_2023_11_10_outgoing_handler_own_future_incoming_response_t ret;
    wasi_http_0_2_0_rc_2023_11_10_outgoing_handler_error_code_t handler_err;
    if (!wasi_http_0_2_0_rc_2023_11_10_outgoing_handler_handle(req, NULL, &ret, &handler_err)) {
        printf("Error sending request\n");
        return 5;
    }

    wasi_http_0_2_0_rc_2023_11_10_types_error_code_t finish_err;
    if (!wasi_http_0_2_0_rc_2023_11_10_types_static_outgoing_body_finish(out_body, NULL, &finish_err)) {
        printf("Failed to finish body\n");
        return 10;
    }

    wasi_http_0_2_0_rc_2023_11_10_types_own_pollable_t poll = wasi_http_0_2_0_rc_2023_11_10_types_method_future_incoming_response_subscribe(wasi_http_0_2_0_rc_2023_11_10_types_borrow_future_incoming_response(ret));
    wasi_io_0_2_0_rc_2023_11_10_poll_method_pollable_block(wasi_io_0_2_0_rc_2023_11_10_poll_borrow_pollable(poll));

    wasi_http_0_2_0_rc_2023_11_10_types_result_result_own_incoming_response_error_code_void_t result;
    if (!wasi_http_0_2_0_rc_2023_11_10_types_method_future_incoming_response_get(wasi_http_0_2_0_rc_2023_11_10_types_borrow_future_incoming_response(ret), &result)) {
        printf("failed to get value for incoming request\n");
        return 1;
    }

    if (result.is_err || result.val.ok.is_err) {
        printf("response is error!\n");
        return 2;
    }

    wasi_io_0_2_0_rc_2023_11_10_poll_pollable_drop_own(poll);
    wasi_http_0_2_0_rc_2023_11_10_types_future_incoming_response_drop_own(ret);

    wasi_http_0_2_0_rc_2023_11_10_types_own_incoming_response_t resp = result.val.ok.val.ok;
    wasi_http_0_2_0_rc_2023_11_10_types_status_code_t code = wasi_http_0_2_0_rc_2023_11_10_types_method_incoming_response_status(wasi_http_0_2_0_rc_2023_11_10_types_borrow_incoming_response(resp));
    response_out->status_code = code;

    wasi_http_0_2_0_rc_2023_11_10_types_own_headers_t header_handle = wasi_http_0_2_0_rc_2023_11_10_types_method_incoming_response_headers(wasi_http_0_2_0_rc_2023_11_10_types_borrow_incoming_response(resp));
    wasi_http_0_2_0_rc_2023_11_10_types_list_tuple2_field_key_field_value_t header_list;
    wasi_http_0_2_0_rc_2023_11_10_types_method_fields_entries(wasi_http_0_2_0_rc_2023_11_10_types_borrow_fields(header_handle), &header_list);

    response_out->headers.headers = malloc(sizeof(header_t) * header_list.len);
    response_out->headers.len = header_list.len;
    for (int i = 0; i < header_list.len; i++) {
        char* name = (response_out->headers.headers[i].name = malloc(header_list.ptr[i].f0.len + 1));
        char* value = (response_out->headers.headers[i].value = malloc(header_list.ptr[i].f1.len + 1));
        strncpy(name, (char *) header_list.ptr[i].f0.ptr, header_list.ptr[i].f0.len);
        name[header_list.ptr[i].f0.len] = 0;
        strncpy(value, (const char*) header_list.ptr[i].f1.ptr, header_list.ptr[i].f1.len);
        value[header_list.ptr[i].f1.len] = 0;
    }
    wasi_http_0_2_0_rc_2023_11_10_types_fields_drop_own(header_handle);

    wasi_http_0_2_0_rc_2023_11_10_types_own_incoming_body_t in_body;
    if (!wasi_http_0_2_0_rc_2023_11_10_types_method_incoming_response_consume(wasi_http_0_2_0_rc_2023_11_10_types_borrow_incoming_response(resp), &in_body)) {
        printf("body is error!\n");
        return 3;
    }
    wasi_io_0_2_0_rc_2023_11_10_streams_own_input_stream_t stream;
    if (!wasi_http_0_2_0_rc_2023_11_10_types_method_incoming_body_stream(wasi_http_0_2_0_rc_2023_11_10_types_borrow_incoming_body(in_body), &stream)) {
        printf("stream is error\n");
        return 9;
    }

    wasi_io_0_2_0_rc_2023_11_10_streams_own_pollable_t stream_poll = wasi_io_0_2_0_rc_2023_11_10_streams_method_input_stream_subscribe(wasi_io_0_2_0_rc_2023_11_10_streams_borrow_input_stream(stream));

    int32_t offset = 0;
    while (offset < response_out->body_max_len) {
        wasi_io_0_2_0_rc_2023_11_10_poll_method_pollable_block(wasi_io_0_2_0_rc_2023_11_10_poll_borrow_pollable(stream_poll));
        wasi_io_0_2_0_rc_2023_11_10_streams_list_u8_t body_res;
        wasi_io_0_2_0_rc_2023_11_10_streams_stream_error_t stream_err;
        if (wasi_io_0_2_0_rc_2023_11_10_streams_method_input_stream_read(wasi_io_0_2_0_rc_2023_11_10_streams_borrow_input_stream(stream), response_out->body_max_len - offset, &body_res, &stream_err)) {
            strncpy(response_out->body + offset, (const char *) body_res.ptr, body_res.len);
            offset = offset + body_res.len;
            response_out->body[offset] = 0;
            
            wasi_io_0_2_0_rc_2023_11_10_streams_list_u8_free(&body_res);
        } else if (stream_err.tag == WASI_IO_0_2_0_RC_2023_11_10_STREAMS_STREAM_ERROR_CLOSED) {
            break;
        } else {
            printf("BODY read is error!\n");
            return 6;
        }
    }
    // wasi_http_0_2_0_rc_2023_11_10_types_outgoing_request_drop_own(req);
    // wasi_io_0_2_0_rc_2023_11_10_streams_input_stream_drop_own(stream);
    wasi_http_0_2_0_rc_2023_11_10_types_incoming_response_drop_own(resp);

    return 0;
}
