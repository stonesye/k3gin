package test

import (
	"k3gin/app/util/trace"
	"testing"
)

func TestNewTraceID(t *testing.T) {
	t.Log(trace.NewTraceID())
}
