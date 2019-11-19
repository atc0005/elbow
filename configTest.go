package main

// https://twitter.com/smazero/status/1196790750660505601
// https://play.golang.com/p/NYoOD_rN03p

import (
	"fmt"
)

type config struct {
	port    *int
	address *string
}

// newConfig takes in basic type parameters and returns a struct containing pointer fields
func newConfig(port int, address string) *config {
	return &config{
		port:    &port, // note we have to take the address of the parameter
		address: &address,
	}
}

// newConfigPointer takes in pointer parameters to initialize the struct
func newConfigPointer(port *int, address *string) *config {
	return &config{
		port:    port,
		address: address,
	}
}

func main() {
	c1 := newConfig(8000, "localhost")

	// print the struct (note the pointer fields so we can't see the values)
	fmt.Println("raw struct", c1)

	// to get the values we have to dereference
	fmt.Println("dereferenced int", *c1.port)
	fmt.Println("dereferenced string", *c1.address)

	// following line won't work as parameters are pointers so require an address not a literal
	// c2 := newConfigPointer(8000, "localhost")

	// nor will this work, as we cannot get the address of a literal directly
	// c2 := newConfigPointer(&8000, &"localhost")

	// we can use the pointer fields from the other struct though
	c2 := newConfigPointer(c1.port, c1.address)
	fmt.Println("dereferenced int on c2", *c2.port)
	fmt.Println("dereferenced string on c2", *c2.address)

	// or we could use variables which we can then get the address of
	port := 9000
	address := "localhost"

	c3 := newConfigPointer(&port, &address)
	fmt.Println("dereferenced int on c3", *c3.port)
	fmt.Println("dereferenced string on c3", *c3.address)
}
