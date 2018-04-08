// Package cas exposes a customizable Content-Adressable Storage API
// with a plugable backend.
//
// By default sha1 is used but others might be used, instead of simply
// hashing the content of the block, a header is also included, this header
// contains:
//
// * type of the content
// * length of the content
//
// Adding both fields makes collisions more complex, since an adversary would
// need to create another payload with a valid structure (defined by type)
// and the same length.
//
// Extra security can added by using different hash options
package cas
