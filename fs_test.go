package store

import (
	"flag"
	"log"
	"os"
	"testing"
)

var (
	tstStorePath string
	st           Store
)

func init() {
	tstStorePath = "/tmp/test-store"
}

func cleanup() {
	if _, err := os.Stat(tstStorePath); !os.IsNotExist(err) {
		os.RemoveAll(tstStorePath)
	}
}

func TestMain(m *testing.M) {
	flag.Parse()
	m.Run()
	//cleanup()
}

func TestCRUD(t *testing.T) {
	st, err := UseFileStore(tstStorePath)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	name := "SuperStore.json"

	// Let Store create a Storage Object out of itself
	err = st.Create(name, st)
	if err != nil {
		t.Error(err)
	}

	var st2 FileStore
	err = st.ReadObject(name, &st2)
	if err != nil {
		t.Error(err)
	}

	if st2.Basepath != st.Basepath {
		log.Printf("st %+v", st)
		log.Printf("st2 %+v", st2)
		t.Error("st2 and st do not equal")

	}

	// Now we are going to try an update
	st2.Comment = "Changed"
	if err = st.Update(name, st2); os.IsNotExist(err) {
		t.Error(err)
	}

	var st3 FileStore
	if err = st.ReadObject(name, &st3); err != nil {
		t.Error(err)
	}
	if st3.Comment != "Changed" {
		t.Errorf("ReadObject expected (Changed) got (%s)", st3.Comment)
	}

	names := st.List()
	if len(names) != 1 {
		t.Errorf("names expected (SuperStore.json) got (%+v)", names)
	}

	// Now test Delete
	if err := st.Delete(name); err != nil {
		t.Error(err)
	}

	if st.Exists(name) {
		t.Error("expected object to have been deleted, but it is still here")
	}
}

// TestNewFileStore
func TestFileStore(t *testing.T) {
	st, err := UseFileStore(tstStorePath)
	if err != nil {
		t.Error(err)
	}

	if st == nil {
		t.Error("TestNewFileStore expected a store got (nil)")
	}

	if _, err := os.Stat(tstStorePath); os.IsNotExist(err) {
		t.Errorf("TestNewFileStore expected %s to exist, but does not", tstStorePath)
	} else if err != nil {
		t.Error(err)
	}
}
