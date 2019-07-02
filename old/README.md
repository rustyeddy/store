# Store ~ Local Object Storage made Simple

Store is a simple Go library to help store application objects on a
local filesystem.  By default, application objects that your store 
to disk are converted to JSON, and unravel when you read them back
in.

## What is Store for?

Store is handy for a couple reaons:

1. Quickly start "persisting" data in a new application without having
   to configure _anything_, no servers, applications, databases or
   containers. 
   
   Just link in the library, tell store what directory to use (store
   will create it if it needs to (and has appropriate permisison, of
   course!))
   
   Then start sticking stuff in, and taking it back out. Simple as
   that!
   
> Just import store into your app and start using it

2. Self contained server / utilities take advantage of Go's ability to
   compile self contained binaries.  

   By adding Store to the application, the application can just start
   using local storage without having to write a bunch of redundant
   scaffold and error checking code.
   
> This will should for in memory filesystems as well, provided it
> looks and acts like a disk based filesystem.

Store can store any type of regular file, images, pdfs, binary blobs
of unknown bit matter.
	 
The following code snippets are how to get started using store quickly.

```go
import "github.com/rustyeddy/store"

func main () {
  // Use a local hidden dir for storage
  st, err  := UseStore(".shoes")
  ifErrorPukeAndDie(err)
  
  // Get a new *Shoe object from the current application
  shoes := NewShoe("hush-puppies")
  
  // Store the shoe: will translate the hush-puppies *Shoes into JSON
  // then write the JSON object as a file named "gumsole.json" in
  // the local directory ".store".
  err := st.StoreObject("gumsole", shoes)
  
  // To get the shoe back, create a *Shoe pointer with space allocated
  // to copy the shoes into from the JSON file
  newShoes := make(Shoes)
  err := st.FetchObject("gumsole", newShoes)
  
  if ! newShoes.Equal(shoes) {
     log.Fatal("I have the wrong Shoes!  Expected (%v) got (%v)", shoes, newShoes)
  }
  
  fmt.Println("Momma gotta new pair a shoes!")
}
```	 

The above example translates the Go _shoes_ object into a JSON string,
then writes that string to the file _gumsole.json_, which is the
object label (name, index, etc) plus the .json representing the saved
file.

All files _should_ have extensions appropriate to the content-type,
store will respect the file extensions.

That is about it.


