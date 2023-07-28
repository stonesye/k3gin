package test

import (
	"k3gin/app/util/conv"
	"testing"
)

func TestParseStringSliceToUint64(t *testing.T) {

	s := make([]string, 2)
	s[0] = "10"
	s[1] = "16"
	s = append(s, "64")
	t.Log(conv.ParseStringSliceToUint64(s))
}
