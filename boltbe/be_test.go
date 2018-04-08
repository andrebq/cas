package boltbe_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/andrebq/cas/boltbe"
)

func TestPutGet(t *testing.T) {
	t.Parallel()
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	be, err := boltbe.New(filepath.Join(dir, "testdb"))
	if err != nil {
		t.Fatal(err)
	}

	w := be.Put([]byte("k"), 1)
	if _, err := w.Write([]byte("content")); err != nil {
		t.Fatal(err)
	}
	if err = w.Close(); err != nil {
		t.Fatal(err)
	}

	r := be.Get([]byte("k"))
	out := make([]byte, 7)
	if _, err := r.Read(out); err != nil {
		t.Fatal(err)
	}
	if string(out) != "content" {
		t.Errorf("expecting %v got %v", "content", string(out))
	}
}
