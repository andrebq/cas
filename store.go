//go:generate mockgen -destination ./mock_io/rwc.go io ReadWriteCloser
//go:generate mockgen -destination ./mock_cas/kv.go -source store.go
package cas

import (
	"bytes"
	"crypto/sha1"
	"hash"
	"io"
	"math"
	"sync"

	"encoding/binary"

	"github.com/pkg/errors"
)

type (
	// Store implements a CAS on top of a KV backend
	Store struct {
		// hash alg pool
		pool sync.Pool

		// kv used to store data
		kv KV
	}

	// KV is the interface used by the store to save cas data
	KV interface {
		// Put returns a writer. The writer MUST accept at least SZ bytes.
		//
		// If there is no space the first call to Write should return an error.
		//
		// When Close is invoked and returns it is expected that the value will be
		// available on KV if, and only if, no error happened during previous write
		// operations.
		//
		// Close is called regardless of Write operations returning errors, all write
		// operations will stop after the first error.
		Put(k []byte, sz int) io.WriteCloser

		// Get returns a reader for the given key k,
		// if k doesn't exist reader should return io.EOF or some other error at
		// the first invocation of Read method.
		//
		// The value is assumed to be static until Close is called.
		//
		// Close will be called regardless of any previous error to write.
		Get(k []byte) io.ReadCloser
	}

	// Option defines one option for the Store
	Option interface {

		// no need to export this
		apply(*Store) error
	}

	optFn  func(*Store) error
	header struct {
		typeLen int
		bodyLen int
	}
)

const (
	maxContentType = byte(math.MaxInt8)
	maxBodyType    = uint32(math.MaxUint32)
)

func (fn optFn) apply(s *Store) error {
	return fn(s)
}

// WithSHA1 makes Store use sha1 as hasher
func WithSHA1() Option {
	return optFn(func(s *Store) error {
		s.pool = sync.Pool{
			New: func() interface{} {
				return sha1.New()
			},
		}
		return nil
	})
}

// WithKV sets the Store to use the provided KV object
func WithKV(kv KV) Option {
	return optFn(func(s *Store) error {
		s.kv = kv
		return nil
	})
}

// NewStore returns a new store with the given options
func NewStore(opts ...Option) (*Store, error) {
	s := &Store{}
	for _, o := range opts {
		err := o.apply(s)
		if err != nil {
			return nil, err
		}
	}

	if s.pool.New == nil {
		return nil, MissingHashFunc
	}

	if s.kv == nil {
		return nil, MissingKVBackend
	}
	return s, nil
}

// Put stores the given typed content into the key value store
// the hash is calculated from is 1 byte type length, 4 bytes content length, type, body
//
// As output you get the key or an error
func (s *Store) Put(contentType, content []byte) ([]byte, error) {
	hdr, err := makeHeader(contentType, content)
	if err != nil {
		return nil, err
	}

	hasher := s.popHash()
	defer s.pushHash(hasher)

	err = writeEntry(hasher, &hdr, contentType, content)
	if err != nil {
		return nil, errors.Wrap(err, "unable to calculate hash key")
	}

	// TODO(andre): maybe use some form of slice caching or something like this?
	key := make([]byte, 0, hasher.Size())
	key = hasher.Sum(key)

	buf := bytes.Buffer{}
	writeEntry(&buf, &hdr, contentType, content)

	w := s.kv.Put(key, hdr.totalSize())
	closer := &pCloser{c: w}
	defer closer.Close()

	err = writeEntry(w, &hdr, contentType, content)
	if err != nil {
		err = closer.Close()
	}
	return key, err
}

func (s *Store) popHash() hash.Hash {
	return s.pool.Get().(hash.Hash)
}

func (s *Store) pushHash(h hash.Hash) {
	s.pool.Put(h)
}

func makeHeader(contentType, content []byte) (header, error) {
	hdr := header{
		typeLen: len(contentType),
		bodyLen: len(content),
	}

	if err := hdr.valid(); err != nil {
		return header{}, err
	}
	return hdr, nil
}

// WriteTo stores on w the binary encoded version of header
func (h header) WriteTo(w io.Writer) (int64, error) {
	var buf [5]byte
	buf[0] = byte(h.typeLen)
	binary.BigEndian.PutUint32(buf[1:], uint32(h.bodyLen))
	sz, err := w.Write(buf[:])
	return int64(sz), err
}

func (h header) totalSize() int {
	return 1 + 4 + h.bodyLen + h.typeLen
}

func (h header) valid() error {
	if uint32(h.bodyLen) > maxBodyType {
		return BodyToBig
	}
	if byte(h.typeLen) > maxContentType {
		return ContentTypeTooBig
	}
	return nil
}

func writeEntry(w io.Writer, hdr *header, contentType, content []byte) (err error) {
	_, err = hdr.WriteTo(w)
	if err != nil {
		err = errors.Wrap(err, "unable to write header")
		return
	}
	_, err = w.Write(contentType)
	if err != nil {
		err = errors.Wrap(err, "unable to write content-type")
		return
	}
	_, err = w.Write(content)
	if err != nil {
		err = errors.Wrap(err, "unable to write content")
		return
	}
	return
}
