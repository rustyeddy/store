# A Simple Object Storage Library

This library is used to Read and Write _Go Objects_ to and from the
filesystem. By default the objects are written as serialized as JSON,
and if they are JSON, they are read and de-serialized into the
corresponding Go object being read by the application.

## Plugins!

This library takes plugins that extend the type and location Store can
used to save objects, plugins may also extend the format(s) objects
can be stored as.

For example, by default Go object will be saved on the local
filesystem as JSON objects.  A plugin may read and write .CSV objects
from _Digital Ocean Spaces_.

## Using Store

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
