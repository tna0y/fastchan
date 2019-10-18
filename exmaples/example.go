package main

import (
	"github.com/tna0y/fastchan"
	"fmt"
)
func main() {
	// Create a new fastchan with buffer size 2
	fc := fastchan.New(2)

	// Put an item
	fc.Put(1)

	// Try to put one more item
	if ok, _ := fc.TryPut(2); ok {
		fmt.Println("Success!")
	}

	// Take two items we just put
	a, _ := fc.Get()
	b, _ := fc.Get()

	// Will print "1 2"
	fmt.Println(a, b)
}
