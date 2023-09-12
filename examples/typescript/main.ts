// @ts-ignore
import { Console } from "as-wasi/assembly";

export function cabi_realloc(a: usize, b: usize, c: usize, len: usize): usize {
  return heap.alloc(len);
}

namespace HTTP {
    export function handle(req: usize, res: usize): void {
        heap.alloc(5);
    }
}

console.log("foo");