package test

import (
	"k3gin/app/util/structure"
	"testing"
)

type s struct {
	Name string
	Age  int
	Sex  int
}

type ts struct {
	Name  string
	Age   int
	Class string
}

func TestCopy(t *testing.T) {
	s := s{
		"testing",
		10,
		1,
	}
	var ts = new(ts)
	err := structure.Copy(s, ts)
	if err != nil {
		t.Errorf("structure copy data : %s ", err.Error())
	}

	t.Log(ts)
}
