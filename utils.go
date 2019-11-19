package main

import "os"

// RemoveDir does what it sounds like it does
func RemoveDir(tpath string) string {
	os.RemoveAll(tpath)
	os.MkdirAll(tpath, 0755)
	return tpath
}
