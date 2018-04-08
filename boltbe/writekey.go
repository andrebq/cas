package boltbe

import (
	"bytes"

	"github.com/boltdb/bolt"
)

type (
	wk struct {
		k      []byte
		db     *bolt.DB
		bucket []byte
		data   bytes.Buffer
	}
)

// Write acculumantes in bytes into a bytes.Buffer to be later flushed when close
// is called
func (w *wk) Write(in []byte) (int, error) {
	return w.data.Write(in)
}

// Close updates the database with the data written to the internal buffer
func (w *wk) Close() error {
	return w.db.Update(func(t *bolt.Tx) error {
		b, err := t.CreateBucketIfNotExists(w.bucket)
		if err != nil {
			return err
		}
		return b.Put(w.k, w.data.Bytes())
	})
}
