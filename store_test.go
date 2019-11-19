package store

import (
	"testing"
)

type ExpectString struct {
	Input  string
	Expect string
}

func TestNameFromPath(t *testing.T) {
	tlist := []ExpectString{
		{"rusty", "rusty"},
		{"/path/to/rusty", "rusty"},
		{"/path/with/slash/", "slash"},
		{"/", ""},
	}
	for _, tst := range tlist {
		name := NameFromPath(tst.Input)
		if name != tst.Expect {
			t.Errorf("NameFromPath %s expected (%s) got (%s)", tst.Input, tst.Expect, name)
		}
	}
}
