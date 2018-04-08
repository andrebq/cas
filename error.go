package cas

type (
	// Error is a static error
	Error string
)

const (
	// ContentTypeTooBig indicates a content-type with a name too large to be stored
	ContentTypeTooBig = Error("content-type is too big. must be less than 255 bytes")

	// BodyToBig indicates a body too big to be processed
	BodyToBig = Error("body is too big, must be less than 4GiB")

	// MissingHashFunc indicates there was no option defined to set a hash function tube used
	MissingHashFunc = Error("missing option to configure hash function to use")

	// MissingKVBackend indicates there was no option defined to set a KV backend
	MissingKVBackend = Error("missing kv backend")
)

// Error implements error interface
func (ce Error) Error() string {
	return string(ce)
}
