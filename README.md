# Fastchan
[![Build Status](https://travis-ci.org/tna0y/Python-random-module-cracker.svg?branch=master)](https://travis-ci.org/tna0y/Python-random-module-cracker)

Go channels, but faster
## When to use

## Install
`go get github.com/tna0y/fastchan`
## How to use
The following example pretty much covers the entire API provided by Fastchan
```go
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
```
## Limitations
* Only buffered channels are supported (for now)
* Buffer size is rounded up to the nearest power of two

## Perfomance
Fastchan may be up to 4.5x times faster

```
BenchmarkFastChan1To1-8   	100341570	        12.1 ns/op
BenchmarkFastChanNTo1-8   	36817144	        32.6 ns/op
BenchmarkFastChan1ToN-8   	19721079	        55.1 ns/op
BenchmarkFastChanNToN-8   	21083403	        56.9 ns/op
```
Equivalent benchmarks for channels
```
BenchmarkChan1To1-8       	21440000	        54.4 ns/op
BenchmarkChanNTo1-8       	19220120	        60.5 ns/op
BenchmarkChan1ToN-8       	12893852	        87.8 ns/op
BenchmarkChanNToN-8       	17704915	        71.5 ns/op
```
## Credits
* Implementation idea [1024cores.net](http://www.1024cores.net/home/lock-free-algorithms/queues/bounded-mpmc-queue)
* Original implementation [Workiva/go-datastructures](https://github.com/Workiva/go-datastructures/blob/master/queue/ring.go)