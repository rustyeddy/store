/*

Package store is a simple CRUD library that applications can use to
read and write Go objects as JSON.  Store can also be used to
store and track other data types, like gif and PDFs.

Every store represents a single level (no subdirectories) storage
container that can be backed by any number of storage continers.

*/
package store

import (
	"log"
	"path/filepath"
	"strings"
	"time"
)

// Store defines how use this package
type Store interface {
	Create(name string, obj interface{}) (interface{}, error)
	ReadObject(name string, obj interface{}) error
	Update(name string, obj interface{})
	Delete(name string) error
}

// Configuration handles all configuration items
type Configuration struct {
	Debug bool
}

var (
	config Configuration
)

func init() {
	config.Debug = false
}

// NameFromPath will extract the name of an object. This is basically just a
// file or directory name.
func NameFromPath(path string) (name string) {
	_, fname := filepath.Split(path)
	dir := filepath.Dir(path)
	if fname == "" && dir != "" {
		_, fname = filepath.Split(dir)
	}

	flen := len(fname) - len(filepath.Ext(fname))
	name = fname[0:flen]
	if config.Debug {
		log.Printf("NameFromPath(%s) returns %s", path, name)
	}
	return name
}

// timeStamp returns a timestamp in a modified RFC3339 format,
// basically remove all colons ':' from filename, since they have a
// specific use with Unix pathnames, hence must be escaped when used
// in a filename.
func timeStamp() string {
	ts := time.Now().UTC().Format(time.RFC3339)
	return strings.Replace(ts, ":", "", -1) // get rid of offesnive colons
}
