package s3

import (
	"fmt"

	"github.com/rustyeddy/stash"
)

type Store struct {
	path  string
	index stash.Index
}

func UseStore(path string) (st *Store, err error) {

	st = &Store{path, make(stash.Index)}

	// 1. Find bucket from GCP - create if necessary

	// 2. Index bucket if it has objects in it

	return st, nil
}

func (st *Store) String() string {
	return fmt.Sprintf("%s objects %d\n", st.path, len(st.Index()))
}

// Path to the storage object e.g. "gcp://storage"
func (st *Store) Path() string {
	return st.path
}

// Name may not be needed?
func (st *Store) Name() string {
	return st.path
}

// Exists returns the number of objects in storage
func (s *Store) Exists(idx string) bool {
	index := s.Index()
	if _, e := index[idx]; e {
		return true
	}
	return false
}

// StoreObject will make sure the object is created and respective
// content is store in the storage
func (st *Store) StoreObject(name string, content interface{}) (obj *stash.Object, err error) {
	panic("todo implement gcp.StoreObject")
	return obj, err
}

// FetchObject will look for the object specified by the Path
// if it exists the object will be returned
func (st *Store) FetchObject(idx string, otype interface{}) (err error) {
	panic("TODO: implement FetchObject")
	return err
}

// DeleteObject removes the object from storage, if it exists
func (st *Store) DeleteObject(idx string) (err error) {
	panic("TODO: implement DeleteObject")
	return err
}

// Index will create an index of objects in storage
func (st *Store) Index() stash.Index {
	panic("TODO: implement Index")
}

// Count returns the number of objects in the store
func (st *Store) Count() int {
	return len(st.index)
}
