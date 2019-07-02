package fs

import (
	"os"
	"testing"
)

// BadStoreTest will ensure that the library reports an error when we
// ask it to use an illegal path for storage.
func TestBadStore(t *testing.T) {
	badpath := "/badpath/dont/be/root/"
	if st, err := UseStore(badpath); err == nil {
		t.Errorf("path should not have been created")
	} else if st != nil {
		t.Errorf("badpath expected (err) got (%v) ", st)
	}
}

func TestCreate(t *testing.T) {

}

func TestRead(t *testing.T) {

}

func TestUpdate(t *testing.T) {

}

func TestDelete(t *testing.T) {

}

// <<<<<<<<<<<<<<<<<<<<<<<<<REMOVE>>>>>>>>>>>>>>>>>>>>>

var (
	testpath string = ".store"
	st       *Store
)

type tKV struct {
	K string
	V string
}

type tKeyVal struct {
	Key string
	Val interface{}
}

// NewStoreTest will test using the newstore
func TestNewStore(t *testing.T) {
	var err error
	var st *Store

	if st, err = UseStore(testpath); err != nil {
		t.Error("expected (store) got error (%v)", err)
	}

	if _, err := os.Stat(st.path); os.IsNotExist(err) {
		t.Error("dir %s does not exist", st.path)
	}
}

// PutStoreTest put some objects in storage
func TestPutStore(t *testing.T) {
	var (
		err error
		st  *Store
		v1  []tKV
		v2  []tKeyVal
	)
	v1 = []tKV{
		{"one", "1"},
		{"two", "2"},
		{"three", "3"},
	}
	v2 = []tKeyVal{
		{"1", 1},
		{"2", 2},
		{"3", 3},
	}

	if st, err = UseStore(testpath); err != nil {
		t.Error(testpath, err)
	}

	var oa, ob *Object
	if _, err = st.StoreObject("idx1", v1); err != nil {
		t.Error("idx1", err)
	}

	if _, err = st.StoreObject("index2", v2); err != nil {
		t.Error("index2", err)
	}

	// lets check that we have two items in the index
	index := st.Index()
	if len(index) != 2 {
		t.Error("index expected (2) got (%d) ", len(index))
	}

	for n, o := range index {
		switch n {
		case "idx1":
			oa = o
		case "index2":
			ob = o
		default:
			t.Errorf("expected index got (%s) ", o.name)
		}
	}

	if oa == nil {
		t.Errorf("expected (idx1) got (%s)", "")
	}
	if ob == nil {
		t.Errorf("expected (idx2) got (%s)", "")
	}
}

func TestExistingStore(t *testing.T) {
	var (
		st  *Store
		err error
	)

	if st, err = UseStore(testpath); err != nil {
		t.Error(testpath, err)
	}

	if len(st.index) != 2 {
		t.Errorf("index len expected (2) got (%d) ", len(st.index))
	}
	found1 := false
	found2 := false

	for n, _ := range st.index {
		switch n {
		case "idx1":
			found1 = true
		case "index2":
			found2 = true
		default:
			t.Errorf("found unexpected index (%d)", n)
		}
	}
	if !found1 || !found2 {
		t.Errorf("expected to find idx1 got (%t) and index2 got (%t)", found1, found2)
	}
}

// TestGetObjects makes sure we can get all the objects we have put
// into the store.
func TestGetObjects(t *testing.T) {
	var (
		st     *Store
		kv     []tKV
		keyval []tKeyVal
		err    error
	)

	if st, err = UseStore(testpath); err != nil {
		t.Error(err)
	}

	if err = st.FetchObject("idx1", &kv); err != nil {
		t.Errorf("expected idx1 got (%v)", err)
	}

	if len(kv) != 3 {
		t.Errorf("expected kv len (3) got (%d) %+v", len(kv), kv)
	}

	for _, v := range kv {
		failed := true
		switch v.K {
		case "1":
			failed = (v.V == "one")
		case "2":
			failed = (v.V == "two")
		case "3":
			failed = (v.V == "three")
		default:
			failed = false
		}
		if failed {
			t.Error("failed", v.K, v.V)
		}
	}

	if err = st.FetchObject("index2", &keyval); err != nil {
		t.Error("expected index2 got (%v)", err)
	}
	if len(keyval) != 3 {
		t.Error("expected keyval len (%d) got (%d)", 2, len(keyval))
	}
}

func TestDeleteObject(t *testing.T) {
	var (
		st  *store.Store
		err error
	)

	if st, err = UseStore(testpath); err != nil {
		t.Error(err)
	}

	if st.Count() != 2 {
		t.Error("expected (2) objects got (%d) ", st.Count())
	}

	idx := "index2"
	if err = st.DeleteObject(idx); err != nil {
		t.Error("failed to delete ", idx)
	}

	if st.Count() != 1 {
		t.Error("count expected (1) got (%d)", st.Count())
	}
}
