package store

import "os"

// RemoveDir will recursively remove an entire directory. Very Dangerous!
// TODO: optionaly create a backup before deleting the entire directory
func RemoveDir(tpath string) string {
	os.RemoveAll(tpath)
	os.MkdirAll(tpath, 0755)
	return tpath
}
