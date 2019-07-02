package store

import (
	"os"
	"testing"
)

type ExpectString struct {
	Input  string
	Expect string
}

// RemoveDir does what it sounds like it does
func RemoveDir(tpath string) string {
	os.RemoveAll(tpath)
	os.MkdirAll(tpath, 0755)
	return tpath
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
