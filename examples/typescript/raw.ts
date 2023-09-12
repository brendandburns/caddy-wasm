export type WasiHandle = i32;
type Char8 = u8;
type Char32 = u32;
export type WasiPtr<T> = usize;
type WasiMutPtr<T> = usize;
type WasiStringBytesPtr = WasiPtr<Char8>;

@external("types", "new-outgoing-request")
export declare function new_outgoing_request(
    method: WasiHandle,
    method_ptr: WasiPtr<Char8>,
    method_len: usize,
    path_ptr: WasiPtr<Char8>,
    path_len: usize,
    query_ptr: WasiPtr<Char8>,
    query_len: usize,
    scheme_is_some: usize,
    scheme: usize,
    scheme_ptr: WasiPtr<Char8>,
    scheme_len: usize,
    authority_ptr: WasiPtr<Char8>,
    authority_len: usize,
    headers: usize,
): WasiHandle;

@external("types", "outgoing-request-write")
export declare function outgoing_request_write(request: WasiHandle, stream_ptr: WasiPtr<WasiHandle>): void;

@external("types", "new-fields")
export declare function new_fields(
    fields_ptr: WasiPtr<usize>,
    fields_len: usize,
): WasiHandle;

@external("default-outgoing-HTTP", "handle")
export declare function handle(
  request: WasiHandle,
  a: usize,
  b: usize,
  c: usize,
  d: usize,
  e: usize,
  f: usize,
  g: usize,
): WasiHandle;

@external("types", "future-incoming-response-get")
export declare function future_incoming_response_get(handle: WasiHandle, ptr: WasiPtr<WasiHandle>): void;

@external("types", "incoming-response-status")
export declare function incoming_response_status(handle: WasiHandle): usize;

@external("types", "incoming-response-headers")
export declare function incoming_response_headers(handle: WasiHandle): WasiHandle;

@external("types", "incoming-response-consume")
export declare function incoming_response_consume(handle: WasiHandle, ptr: WasiPtr<u8>): void;

@external("streams", "read")
export declare function streams_read(handle: WasiHandle, len: i64, ptr: WasiPtr<Char8>): void;

@external("streams", "write")
export declare function streams_write(handle: WasiHandle, ptr: WasiPtr<Char8>, len: usize, result: WasiPtr<WasiHandle>): void;


@unmanaged
export class WasiString {
    ptr: WasiStringBytesPtr;
    length: usize;

    constructor(str: string) {
        let wasiString = String.UTF8.encode(str, false);
        // @ts-ignore: cast
        this.ptr = changetype<WasiStringBytesPtr>(wasiString);
        this.length = wasiString.byteLength;
    }

    toString(): string {
        let tmp = new ArrayBuffer(this.length as u32);
        memory.copy(changetype<usize>(tmp), this.ptr, this.length);
        return String.UTF8.decode(tmp);
    }
}