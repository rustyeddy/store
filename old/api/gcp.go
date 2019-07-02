package gcp

import (
	"context"

	log "github.com/rustyeddy/logrus"
	"github.com/rustyeddy/stash"
	"google.golang.org/api/iterator"

	"cloud.google.com/go/storage"
)

var (
	cli     *storage.Client
	buckets []string // list of bucket names
)

func init() {
}

// GetClient for gcp storage
func GetClient() *storage.Client {
	var err error
	ctx := context.Background()
	if cli == nil {
		if cli, err = storage.NewClient(ctx); err != nil {
			log.Fatal("failed to get GCP client")
		}
	}
	return cli
}

// GetOrCreateBucket will return a bucket *pointer to the the bucket
// with the given name.  The bucket will be created if it does not
// already exist.
func (st *Store) GetOrCreateBucket(ctx context.Context, path string) *storage.BucketHandle {
	cli := GetClient()
	st.bucket = cli.Bucket(path)
	attrs, err := st.bucket.Attrs(ctx)
	if err != nil {
		if err := st.bucket.Create(ctx, st.projid, nil); err != nil {
			log.Fatal("failed to create bucket", err)
		}
	}
	log.Fatalf("%+v", attrs)
	return st.bucket
}

// Buckets returns a list of bucket names for the current project
func (st *Store) Buckets() ([]string, error) {
	ctx := context.Background()

	it := st.client.Buckets(ctx, st.projid)
	for {
		battrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, battrs.Name)
	}
	return buckets, nil
}

// Index returns a map of ObjectHandles indexed
func (st *Store) Index() (index stash.Index) {
	ctx := context.Background()
	bkt := st.bucket
	it := bkt.Objects(ctx, nil)
	for {
		oattrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("failed ", err)
			return nil
		}
		st.index[oattrs.Name] = oattrs
	}
	return st.index
}

func (st *Store) Names() []string {
	ctx := context.Background()

	bkt := st.bucket
	var objs []string
	it := bkt.Objects(ctx, nil)
	for {
		oattrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("failed ", err)
			return nil
		}
		objs = append(objs, oattrs.Name)
	}
	return objs
}
