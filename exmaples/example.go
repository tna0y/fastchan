package main

import (
	"fmt"
	"github.com/tna0y/fastchan"
)

func main() {
	// Create a new fastchan with buffer size 2
	fc := fastchan.New(2)

	// Put an item
	fc.Put(1)

	// Try to put one more item
	if ok := fc.TryPut(2); ok {
		fmt.Println("Success!")
	}

	// Take two items we just put
	a := fc.Get()
	b := fc.Get()

	// Will print "1 2"
	fmt.Println(a, b)
}
