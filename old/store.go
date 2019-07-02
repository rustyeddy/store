package store

/*
	Store is a place to store things
*/

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// ====================================================================
//                        Store
// ====================================================================

// Store is the main structure of this package.  It has a name and
// maintains a path, ObjectIndex and some house keeping private fields
type Store struct {
	Path    string // basedir of this store
	Name    string // the name of the store provider
	Created string
	index

	Fetched int64
	Stored  int64
	Errored int64
	Indexed int64
}

// UseStore
// ====================================================================

// UseStore creates and returns a new storage container.  If a dir
// already exists, that directory will be used.  If the directory
// does not exist it will be created.
func UseStore(path string) (s *Store, err error) {
	path = filepath.Clean(path)
	s = &Store{
		Path:    path,
		Name:    NameFromPath(path),
		Created: timeStamp(),
		index:   make(index),
	}

	// Determine if we are using an existing directory or need
	// to create a new one.
	if _, err = os.Stat(path); os.IsNotExist(err) {
		// create the path that did not previously exist
		if err := os.MkdirAll(path, 0755); err != nil {
			return nil, fmt.Errorf("mkdir %s: %v", path, err)
		}
	}
	// We already have the store get an index
	s.buildIndex()

	log.Debugln(s.String())
	return s, nil
}

// String a simple summary of our store
func (s *Store) String() string {
	return fmt.Sprintf("store %q, path %q, objects %d",
		s.Name, s.Path, len(s.Index()))
}

// Indexing
// ====================================================================

// Index gives us a map indexed by filenames with pointers to
// the corresponding object.  The index needs to be rebuilt as
// the result of any change to the underlying filesystem.  Likewise
// changes to index or any in memory storage will need to be
// flushed to the underlying storage.
func (s *Store) Index() index {
	if s.index == nil {
		s.buildIndex()
	}
	return s.index
}

// FilterNames just makes sure index is built and calls the
// respective index function
func (s *Store) FilterNames(f func(name string) string) (names []string, objs []*Object) {
	i := s.Index()
	names, objs = i.FilterNames(f)
	return names, objs
}

// timeStamp returns a timestamp in a modified RFC3339 format,
// basically remove all colons ':' from filename, since they have a
// specific use with Unix pathnames, hence must be escaped when used
// in a filename.
func timeStamp() string {
	ts := time.Now().UTC().Format(time.RFC3339)
	return strings.Replace(ts, ":", "", -1) // get rid of offesnive colons
}

// =======================================================================
// Index returns a map of item names and full paths
// =======================================================================

// Index will scan the store directory for objects (files) creating a
// map of pointers to the Objects indexed by the object name (file
// name less the path and extension)
func (s *Store) buildIndex() index {
	// Now build the index if we don't have one
	s.Path = filepath.Clean(s.Path) // Cleanse our path
	pattern := s.Path + "/*"
	s.indexPaths(pattern)
	return s.index
}

// indexPaths will create a map of *File created from fullpaths indexed by
// the filename (less extension).
func (s *Store) indexPaths(pattern string) (err error) {
	var (
		paths []string
	)

	if paths, err = filepath.Glob(pattern); err != nil || paths == nil {
		return fmt.Errorf("no files to index %s %v", pattern, err)
	}

	// Create room in the index for the paths
	if s.index == nil {
		s.index = make(index, len(paths))
	}

	// for the range of paths
	for _, p := range paths {
		var (
			fi  os.FileInfo
			err error
		)

		// We only want to index regular files, Lstat will help use determine
		if fi, err = os.Lstat(p); err != nil {
			log.Warningln("Lstat error ", p, err) // TODO: Append to a buffer
			continue
		}

		// We want to log this incase a time comes later if we do
		// care about non-regular files
		if !fi.Mode().IsRegular() {

			// Should we complain about directories?
			log.Debugf("  ignore (non regular file) %s", p)
			continue
		}

		var obj *Object
		if obj, err = ObjectFromPath(p); err != nil {
			log.Errorln(p, err)
		}

		obj.AddInfo(fi)
		obj.Store = s

		// attach the object to the index
		s.Set(obj.Name, obj)
	}
	return nil
}

// Count returns the number of items in Store
func (s *Store) Count() int {
	return len(s.Index())
}

// ==============================================================================

func (s *Store) Shutdown() {
	// Flush in memory data
	// Shutdown connections
}

func (s *Store) 
		fstats := getFileStats(fi)
