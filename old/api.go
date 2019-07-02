/*

Store is a classic CRUD library that stores generic Go
objects as JSON files.

*/
package store

// CRUD Interface defines the base operations we expect from our backer
type Store interface {
	Create(name string, obj interface{}) (Object, error) // Create an object
	Read(name string) (Object, error)                    // Read an Object
	Update(name string, obj interface{}) error           // Update an Object
	Delete(name string) error                            // Delete the object

	Index() map[string]Object
	Size() int
}

// Object defines the holding structure for the generic data that makes
// up a particular object
type Object interface {
	Type() string
	Size() string
	io.Read
	io.Write

	//Read() (buf []byte, obj int)
	//Write(buf []byte) int
}
