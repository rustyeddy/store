# A Simple Object Storage Library

This library is used to Read and Write _Go Objects_ to and from the
filesystem (persistance). By default the objects are written as serialized as JSON,
and if they are JSON, they are read and de-serialized into the
corresponding Go object being read by the application.

## Using Store

Store provides the following very simple interface:

```go
type Store interface {
	Create(name string, obj interface{}) (interface{}, error)
	ReadObject(name string, obj interface{}) error
	Update(name string, obj interface{})
	Delete(name string) error
}
```

Anything that provides the above interface can be used as a backend.  For example, we can add AWS S3, 
Digital Ocean storage and Google Cloud Project.

```go
package main 

import "github.com:rustyeddy/store"

func crud_example(tstStorePath string) {
	st, err := UseFileStore(tstStorePath)
	if err != nil {
		log.Fatal(err)
	}
	name := "SuperStore.json"

	// Let Store create a Storage Object out of itself
	err = st.Create(name, st)
	if err != nil {
		log.Fatal(err)
	}
	
	// Now read that object 
	var st2 FileStore
	err = st.ReadObject(name, &st2)
	if err != nil {
		log.Fatal(err)
	}

	// Now we are going to try an update
	st2.Comment = "Changed"
	if err = st.Update(name, st2); os.IsNotExist(err) {
		log.Fatal(err)
	}

	var st3 FileStore
	if err = st.ReadObject(name, &st3); err != nil {
		log.Fatal(err)
	}
	if st3.Comment != "Changed" {
		log.Fatalf("ReadObject expected (Changed) got (%s)", st3.Comment)
	}

	// Now test Delete
	if err := st.Delete(name); err != nil {
		;log.Fatal(err)
	}

	if st.Exists(name) {
		log.Fatal("expected object to have been deleted, but it is still here")
	}
}

int main() {
    crud_example(".")
}

```

That is all it takes to get started.
