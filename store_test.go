package cas_test

import (
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/andrebq/cas"
	"github.com/andrebq/cas/boltbe"
	"github.com/andrebq/cas/mock_cas"
	"github.com/andrebq/cas/mock_io"
	"github.com/golang/mock/gomock"
)

func TestGet(t *testing.T) {
	// using boltbe to simplify tests

	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	kv, err := boltbe.New(filepath.Join(tmpdir, "be"))
	if err != nil {
		t.Fatal(err)
	}

	store, err := cas.NewStore(
		cas.WithSHA1(),
		cas.WithKV(kv))
	if err != nil {
		t.Fatal(err)
	}

	type testEntry struct {
		contentType []byte
		content     []byte
		expectedErr error
		out         []byte
	}

	tests := []testEntry{
		{
			contentType: []byte("t"),
			content:     []byte("hello"),
			expectedErr: nil,
			out:         nil,
		},
		{
			contentType: []byte("t"),
			content:     []byte("hello"),
			expectedErr: nil,
			out:         make([]byte, 6),
		},
	}

	for _, test := range tests {
		k, err := store.Put(test.contentType, test.content)
		if err != nil {
			t.Fatal(err)
		}
		actualContentType, actualContent, err := store.Get(test.out, k)
		if err != test.expectedErr {
			t.Fatalf("Expecting error %v got %v", test.expectedErr, err)
		} else if !reflect.DeepEqual(actualContent, test.content) {
			t.Fatalf("Expecing %v got %v for content", string(test.content), string(actualContent))
		} else if !reflect.DeepEqual(actualContentType, test.contentType) {
			t.Fatalf("Expecing %v got %v for content-type", string(test.contentType), string(actualContentType))
		}
	}
}

func TestPut(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rwc := mock_io.NewMockReadWriteCloser(ctrl)
	rwc.EXPECT().Write(gomock.Any()).Times(3 /* hdr write + content type + content */).Return(0, nil)
	rwc.EXPECT().Close().Return(nil).Times(1)

	kv := mock_cas.NewMockKV(ctrl)
	kv.EXPECT().Put(gomock.Any(), gomock.Any()).Return(rwc).Times(1)

	store, err := cas.NewStore(
		cas.WithSHA1(),
		cas.WithKV(kv))

	if err != nil {
		t.Fatal(err)
	}
	type testEntry struct {
		contentType []byte
		content     []byte
		expectedErr error
		expectedKey []byte
	}
	updateKey := func(te *testEntry) {
		h := sha1.New()
		binary.Write(h, binary.BigEndian, byte(len(te.contentType)))
		binary.Write(h, binary.BigEndian, uint32(len(te.content)))
		h.Write(te.contentType)
		h.Write(te.content)
		te.expectedKey = h.Sum(nil)
	}

	tests := []testEntry{
		{
			contentType: []byte("t"),
			content:     []byte("hello"),
			expectedErr: nil,
		},
	}
	for _, test := range tests {
		if test.expectedErr == nil {
			updateKey(&test)
		}
		actualKey, err := store.Put(test.contentType, test.content)
		if test.expectedErr != nil {
			if err == nil ||
				test.expectedErr.Error() != err.Error() {
				t.Errorf("For %v/%v should get error %v got %v",
					string(test.contentType),
					string(test.content),
					test.expectedErr,
					err)
				continue
			}
		} else if !reflect.DeepEqual(actualKey, test.expectedKey) {
			t.Errorf("For %v/%v key should be %v got %v",
				string(test.contentType),
				string(test.content),
				hex.EncodeToString(test.expectedKey),
				hex.EncodeToString(actualKey))
		} else if err != nil {
			t.Error(err)
		}
	}
}
