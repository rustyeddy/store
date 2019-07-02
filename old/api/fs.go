package fs

/*
	Store is a place to store things
*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/rustyeddy/logrus"
	"github.com/rustyeddy/stash"
)

// Filestore implements the storage interface with the local filesystem
type Filestore struct {
	Basepath string // No proto means local filesystem
	Objects  map[string]Object
}

// ====================================================================
//                        Store
// ====================================================================

// UseStore creates and returns a new storage container.  If a dir
// already exists, that directory will be used.  If the directory
// does not exist it will be created.
func UseStore(path string) (stash.Storage, error) {
	path = filepath.Clean(path)
	s := &Store{
		path:    path,
		name:    nameFromPath(path),
		created: TimeStamp(),
	}

	// TODO - XXX - Add a permission check.
	// Determine if we are using an existing directory or need
	// to create a new one.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// path did not exist, but it will now if we can help it
		if err := os.MkdirAll(path, 0755); err != nil {
			return nil, fmt.Errorf("mkdir %s: %v", path, err)
		}
	}
	s.Index() // Index the store for first time
	return s, nil
}

// String a simple summary of our store
func (s *Store) String() string {
	return fmt.Sprintf("name: %s path %s, object count: %d",
		s.name, s.path, len(s.index))
}

// Path returns the path in file system
func (s *Store) Path() string {
	return s.path
}

// Name returns the name
func (s *Store) Name() string {
	return s.name
}

// pathFromName derives the path from the name. It combines the
// store Path with the filename.
func (s *Store) pathFromName(name string) string {
	return s.path + "/" + name
}

// TimeStamp returns a timestamp in a modified RFC3339
// format, basically remove all colons ':' from filename, since
// they have a specific use with Unix pathnames, hence must be
// escaped when used in a filename.
func TimeStamp() string {
	ts := time.Now().UTC().Format(time.RFC3339)
	return strings.Replace(ts, ":", "", -1) // get rid of offesnive colons
}

// Exists will tell us if the storage object exists
func (s *Store) Exists(idx string) bool {
	i := s.Index()
	if _, e := i[idx]; e {
		return true
	}
	return false
}

/*
```````````````````````````````````````````````````````````````````````

                     Store and Fetch

    These functions take an empty Go interface, serialize it into a
    the specified format, e.g. JSON, YML, etc. then write the file
    to the underlying container.


```````````````````````````````````````````````````````````````````````
*/

// StoreObject will serialize this JSON Go object and store it
// in the file system.  This command is destructive, any existing
// data will be over written.
func (s *Store) StoreObject(name string, data interface{}) (obj *stash.Object, err error) {
	stobj := data
	jbuf, err := json.Marshal(stobj) // JSONify the "object" param
	if err != nil {
		return nil, fmt.Errorf("JSON marshal failed %s -> %v", name, err)
	}

	// log.Debug("  storing data :", string(jbuf[0:40]))
	obj = stash.ObjectFromBytes(jbuf) // obj will not be nil
	obj.Storage = s                   // back pointer to store
	obj.Name = name
	obj.Path = filepath.Clean(s.pathFromName(name)) + ".json"

	err = ioutil.WriteFile(obj.Path, jbuf, 0644)
	if err != nil {
		return nil, fmt.Errorf("  Store.Object write failed %s -> %v", obj.Path, err)
	}

	// Index this object. This is destructive, it could be overwriting
	// an existing object.  Index should never return nil, we'll ASSUME
	// it will not be nil (it can be empty but must exist).
	if s.index == nil {
		s.index = make(stash.Index)
	}
	s.index[name] = obj
	return obj, nil
}

// FetchObject returns the name object from the filesystem, the content
// is decoded if desired (e.g. JSON), and a Go object is returned. Nil
// is returned if thier is a problem, such as no object existing.
func (s *Store) FetchObject(name string, otype interface{}) error {
	var (
		obj *stash.Object
		ex  bool
		err error
	)
	o, ex := s.index[name]
	if ex == false {
		return fmt.Errorf("Fetch Object does NOT EXIST")
	}
	obj = o.(*stash.Object)

	// If our buffer is nil, we will need to fetch the data from the store.
	if obj.Buffer == nil {
		obj.Buffer, err = ioutil.ReadFile(obj.Path)
		if err != nil {
			return fmt.Errorf("  FetchObject failed reading %s -> %v\n", obj.Path, err)
		}
	}
	log.Debugf("  ++ found %d bytes from %s ", len(obj.Buffer), obj.Path)

	// Determine the content type we are dealing with
	ext := filepath.Ext(obj.Path)
	if ext != "" {
		if obj.ContentType = mime.TypeByExtension(ext); obj.ContentType == "" {
			obj.ContentType = http.DetectContentType(obj.Buffer)
		}
	}

	if obj.ContentType == "application/json" {
		if err := json.Unmarshal(obj.Buffer, otype); err != nil {
			return fmt.Errorf("%s: %v", name, err)
		}
	}
	return nil
}

// DeleteObject does just that, it removes the object from the store.
// Meaning it removes the object from the disk
func (s *Store) DeleteObject(name string) error {
	var (
		obj *stash.Object
	)
	o, e := s.index[name]
	if !e {
		return fmt.Errorf("%s NOT FOUND", name)
	}

	obj = o.(*stash.Object)

	// The object must be removed from the filesystem first ...
	if obj.Path == "" {
		return fmt.Errorf("path is nil, should never happen %s", name)
	}
	if err := os.Remove(obj.Path); err != nil {
		return fmt.Errorf("Remove path %s error %v", obj.Path, err)
	}

	// Now remove form the index.
	delete(s.index, name)
	return nil
}

// =======================================================================
// Index returns a map of item names and full paths
// =======================================================================

// Index will scan the store directory for objects (files) creating a
// map of pointers to the Objects indexed by the object name (file
// name less the path and extension)
func (s *Store) Index() stash.Index {
	// Now build the index if we don't have one
	s.path = filepath.Clean(s.path) // Cleanse our path
	pattern := s.path + "/*"
	s.indexPaths(pattern)
	return s.index
}

// indexPaths will create a map of *File created from fullpaths indexed by
// the filename (less extension).
func (s *Store) indexPaths(pattern string) (stash.Index, error) {
	var (
		paths []string
		err   error
	)

	if paths, err = filepath.Glob(pattern); err != nil || paths == nil {
		return nil, fmt.Errorf("no files to index %s %v", pattern, err)
	}

	// Create room in the index for the paths
	if s.index == nil {
		s.index = make(stash.Index, len(paths))
	}

	// for the range of paths
	for _, p := range paths {
		var (
			fi  os.FileInfo
			err error
		)

		// We only want to index regular files, Lstat will help use determine
		if fi, err = os.Lstat(p); err != nil {
			log.Warning("Lstat %s: error %v: ", p, err) // TODO: Append to a buffer
			continue
		}
		if !fi.Mode().IsRegular() {
			log.Infof("  ignore (non regular file) %s", p)
			continue
		}

		var obj *stash.Object
		if obj, err = stash.ObjectFromPath(p); err != nil {
			log.Errorln(p, err)
		}
		obj.Storage = s

		// attach the object to the index
		s.index[obj.Path] = obj
	}
	return s.index, nil
}

// Count returns the number of items in Store
func (s *Store) Count() int {
	return len(s.index)
}
