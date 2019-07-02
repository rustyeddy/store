package gcp

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/storage"
	log "github.com/rustyeddy/logrus"
	"github.com/rustyeddy/stash"
)

type Store struct {
	path  string
	index stash.Index

	// GCP specific stuff
	bucket *storage.BucketHandle
	client *storage.Client
	projid string
}

func UseStore(path string) (*Store, error) {

	st := &Store{path: path, index: make(stash.Index)}

	ctx := context.WithValue(context.Background(), "path", path)
	if st.bucket = st.GetOrCreateBucket(ctx, path); st.bucket == nil {
		return st, errors.New("Bucket Not Found " + path)
	}
	if index := st.Index(); index == nil {
		return st, errors.New("failed to index bucket objects")
	}
	return st, nil
}

// client returns a GCP storage client
func (st *Store) Client() *storage.Client {
	if st.client == nil {
		if st.client = GetClient(); st.client == nil {
			log.Fatal("failed to get GCP client")
		}
	}
	return st.client
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
	panic("todo implement gcp.Exists()")
}

// StoreObject will make sure the object is created and respective
// content is store in the storage
func (st *Store) StoreObject(name string, content interface{}) (obj *stash.Object, err error) {
	panic("todo implement gcp.StoreObject")
	return obj, err
}

// FetchObject will look for the object specified by the Path
// if it exists the object will be returned
func (s *Store) FetchObject(idx string, otype interface{}) (err error) {
	panic("TODO: implement FetchObject")
	return err
}

// DeleteObject removes the object from storage, if it exists
func (s *Store) DeleteObject(idx string) (err error) {
	panic("TODO: implement DeleteObject")
	return err
}

// Count returns the number of objects in storage
func (s *Store) Count() int {
	return len(s.index)
}
