package test

import (
	"k3gin/app/util/uuid"
	"testing"
)

func TestUUID(t *testing.T) {
	t.Run("NewUUID", func(t *testing.T) {
		t.Log(uuid.NewUUID())
	})

	t.Run("MustString", func(t *testing.T) {
		t.Log(uuid.MustString())
	})
}
