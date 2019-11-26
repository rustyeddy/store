package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

// FileStore is a container that satisfies the Store interface
type FileStore struct {
	Name     string // Name of the File store
	Basepath string // Path to storage directory
	Comment  string // Some useful comment
	Names    []string
}

// UseFileStore creates and returns a new storage container.  If a dir
// already exists, that directory will be used.  If the directory does
// not exist it will be created.
func UseFileStore(path string) (fs *FileStore, err error) {
	path = filepath.Clean(path)
	fs = &FileStore{
		Basepath: path,
		Name:     NameFromPath(path),
		Comment:  "Hello, World!",
	}

	// TODO - Add a permission check, Determine if we are using
	// an existing directory or need to create a new one.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// path did not exist, but it will now if we can help it
		if err := os.MkdirAll(path, 0755); err != nil {
			return nil, fmt.Errorf("mkdir %s: %v", path, err)
		}
	}
	return fs, nil
}

// Path will tell us of the location on the filesystem
func (fs *FileStore) Path() string {
	return fs.Basepath
}

// Create takes the provided Go object, converts it to a JSON string
// (or mimetype), then stores it under the given object name.
func (fs *FileStore) Create(name string, gobj interface{}) (err error) {
	if fs.Exists(name) {
		return fmt.Errorf("File Object Exists")
	}

	// turn the interface into and object then JSONify it
	stobj := gobj

	jbuf, err := json.MarshalIndent(stobj, "  ", "  ")
	if err != nil {
		return err
	}

	path := fs.Basepath + "/" + name

	// Now write to the file
	err = ioutil.WriteFile(path, jbuf, 0644)
	if err != nil {
		return fmt.Errorf("  Store.Object write failed %s -> %v",
			path, err)
	}
	return err
}

// ReadObject will read the object name and if tobj has been supplied it will
// be returned populated or an error will have resulted.
func (fs *FileStore) ReadObject(name string, gobj interface{}) (err error) {
	path := fs.Basepath + "/" + name
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	contentType := ""
	ext := filepath.Ext(name)
	if ext != "" {
		ext = ext[1:]
		if ext == "json" {
			contentType = "application/json"
		} else if contentType = mime.TypeByExtension(ext); contentType == "" {
			contentType = http.DetectContentType(buf)
		}
	}

	// If that object is json then unmarshal this json
	if contentType == "application/json" {
		if err := json.Unmarshal(buf, &gobj); err != nil {
			return fmt.Errorf("Read: %s ~> %v", name, err)
		}
	}
	return err
}

// Update will replace the contents of the named object with the whole
// of the new object.
func (fs *FileStore) Update(name string, gobj interface{}) (err error) {
	// Update is opposite of create. Update errors if the file does not exist
	if !fs.Exists(fs.Basepath + "/" + name) {
		return fmt.Errorf("Update: File Does Exists")
	}

	// XXX ~ We are assuming this is going to be JSON Fix this!
	stobj := gobj // turn an interface into an object

	// JSONify the "object" param
	jbuf, err := json.MarshalIndent(stobj, "  ", "  ")
	if err != nil {
		return err
	}

	path := fs.Basepath + "/" + name

	// Now write to the file
	err = ioutil.WriteFile(path, jbuf, 0644)
	if err != nil {
		return fmt.Errorf("  Store.Object write failed %s -> %v", path, err)
	}
	return err
}

// Delete handles the removal of the specified file
func (fs *FileStore) Delete(name string) (err error) {
	err = os.RemoveAll(fs.Basepath + "/" + name)
	return err
}

// Exists will tell of if the named storage file exists
func (fs *FileStore) Exists(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false
	}
	return true
}

// Save will create the file if it does not already exist, if the file
// does exist it will be updated.
func (fs *FileStore) Save(name string, gobj interface{}) (err error) {
	if fs.Exists(name) {
		err = fs.Update(name, gobj)
	} else {
		err = fs.Create(name, gobj)
	}
	return err
}

// List will provide a list of all elements in storage
func (fs *FileStore) List() (names []string) {

	if fs.Names != nil {
		return fs.Names
	}

	files, err := ioutil.ReadDir(fs.Basepath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fs.Items = append(fs.Names, file.Name())
	}
	return fs.Names
}
