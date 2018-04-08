package cas_test

import (
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/andrebq/cas"
	"github.com/andrebq/cas/mock_cas"
	"github.com/andrebq/cas/mock_io"
	"github.com/golang/mock/gomock"
)

func TestPutGet(t *testing.T) {
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
