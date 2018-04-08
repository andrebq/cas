package boltbe

import "io"

type (
	notFoundReader struct{}
)

// Read implements io.Read but always returns unexpected EOF
func (notFoundReader) Read(out []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

// Close implements io.Close and always returns nil
func (notFoundReader) Close() error {
	return nil
}
