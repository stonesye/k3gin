package test

import (
	"github.com/go-playground/assert/v2"
	"k3gin/app/util/hash"
	"testing"
)

func TestMD5(t *testing.T) {
	hashMD5 := "b5f6e4ead14cbb0ab1434fe25b8f4f58"
	assert.Equal(t, hashMD5, hash.MD5([]byte("1234565789")))
}

func TestMD5String(t *testing.T) {
	t.Run("md5String", func(t *testing.T) {
		hashMD5 := "b5f6e4ead14cbb0ab1434fe25b8f4f58"
		assert.Equal(t, hashMD5, hash.MD5String("1234565789"))
	})
}

func TestSHA1String(t *testing.T) {
	assert.Equal(t, "f7c3bc1d808e04732adf679965ccc34ca7ae3441", hash.SHA1String("123456789"))
}
