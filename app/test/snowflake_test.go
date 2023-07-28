package test

import (
	"k3gin/app/util/snowflake"
	"testing"
)

func TestMustID(t *testing.T) {
	snowflake.Init()
	t.Log(snowflake.MustID())
}
