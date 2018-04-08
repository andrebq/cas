package boltbe

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/boltdb/bolt"
)

type (
	// KV implements cas.KV interface using boltdb as backend
	KV struct {
		db *bolt.DB
	}
)

var (
	theBucket = []byte("content")
)

// New returns a new KV stored at dbfile location
func New(dbfile string) (*KV, error) {
	db, err := bolt.Open(dbfile, 0766, nil)
	if err != nil {
		return nil, err
	}
	return &KV{db: db}, nil
}

// Put implements cas.KV.Put
func (kv *KV) Put(k []byte, sz int) io.WriteCloser {
	return &wk{
		bucket: theBucket,
		db:     kv.db,
		k:      k,
	}
}

// Get implements cas.KV.Get
func (kv *KV) Get(k []byte) io.ReadCloser {
	var buf []byte
	kv.db.View(func(t *bolt.Tx) error {
		b := t.Bucket(theBucket)
		buf = b.Get(k)
		return nil
	})
	if len(buf) == 0 {
		return notFoundReader{}
	}
	return ioutil.NopCloser(bytes.NewBuffer(buf))
}
