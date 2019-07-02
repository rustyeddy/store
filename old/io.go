package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

/* ====================================================================

                     Store and Fetch

    These functions take an empty Go interface, serialize it into a
    the specified format, e.g. JSON, YML, etc. then write the file
    to the underlying container.

===================================================================== */

// StoreObject accepts any Go structure, serialize it as a JSON object.
// The JSON object will then be written to disk encapsulated as an "Object".
// The Object contains some meta data about the original object, including
// it's type to help the Application be able to deserialize and use the
// Object with out implicit knowledge of the objects structure.
func (s *Store) StoreObject(name string, data interface{}) (obj *Object, err error) {
	defer func() {
		if err == nil {
			s.Stored++
		} else {
			s.Errored++
		}
	}()

	// Do NOT allow '/' characters in string
	if strings.Index(name, "/") > -1 {
		return nil, l.IndexError(name, "illegal char '/' used for index")
	}

	stobj := data                                      // turn an interface into an object
	jbuf, err := json.MarshalIndent(stobj, "  ", "  ") // JSONify the "object" param
	if err != nil {
		return nil, l.JSONError(name, err.Error())
	}

	// log.Debug("  storing data :", string(jbuf[0:40]))
	obj = ObjectFromBytes(jbuf) // obj will not be nil
	if obj == nil {
		return nil, errors.New("StoreObject failed translating bytes")
	}
	obj.Store = s // back pointer to store
	obj.Name = name
	obj.Path = s.Path + "/" + name + ".json"

	// Now write to the file
	err = ioutil.WriteFile(obj.Path, jbuf, 0644)
	if err != nil {
		return nil, fmt.Errorf("  Store.Object write failed %s -> %v", obj.Path, err)
	}
	s.Set(name, obj)
	return obj, nil
}

// Fetch Object ~ cRud
// ====================================================================

// FetchObject unmarshal the contents of the file using the type
// template passed in by otype.  The original Go object will be
// is decoded if desired (e.g. JSON), and a Go object is returned. Nil
// is returned if thier is a problem, such as no object existing.
func (s *Store) FetchObject(name string, otype interface{}) (obj *Object, err error) {
	defer func() {
		if err == nil {
			s.Fetched++
		} else {
			s.Errored++
		}
	}()

	if obj = s.Get(name); obj == nil {
		return nil, l.FetchError(" does NOT EXIST ", name)
	}

	// If our buffer is nil, we will need to fetch the data from the store.
	if obj.Buffer == nil {
		obj.Buffer, err = ioutil.ReadFile(obj.Path)
		if err != nil {
			return nil, l.FetchError(obj.Path, err.Error())
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
			return nil, fmt.Errorf("store.FetchObject: %s ~> %v", name, err)
		}
	}
	return obj, nil
}

// DeleteObject ~ Delete
// ====================================================================

// DeleteObject does just that, it removes the object from the store.
// Meaning it removes the object from the disk
func (s *Store) DeleteObject(name string) error {
	var (
		obj *Object
	)
	if obj = s.Get(name); obj == nil {
		return fmt.Errorf("%s NOT FOUND", name)
	}

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
