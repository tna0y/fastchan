# Fastchan
[![Build Status](https://travis-ci.org/tna0y/Python-random-module-cracker.svg?branch=master)](https://travis-ci.org/tna0y/Python-random-module-cracker)

Fastchan is a data structure that aims to create a __faster__ version of Go channels
that preserves all properties and thus may be swapped in easily. 

The aim is to achieve superior performance through the use of lock-free data structures instead of mutexes.

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
    if ok := fc.TryPut(2); ok {
        fmt.Println("Success!")
    }
    
    // Take two items we just put
    a := fc.Get()
    b := fc.Get()
    
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
BenchmarkFastChan1To1-8            	84515737	        13.3 ns/op
BenchmarkFastChanNTo1-8            	33871244	        35.5 ns/op
BenchmarkFastChan1ToN-8            	21209842	        60.8 ns/op
BenchmarkFastChanNToN-8            	20933433	        56.8 ns/op
BenchmarkFastChanBufferedRead-8    	91768263	        13.0 ns/op
BenchmarkFastChanBufferedWrite-8   	89748022	        13.3 ns/op
```
Equivalent benchmarks for channels
```
BenchmarkChan1To1-8                	20599215	        56.1 ns/op
BenchmarkChanNTo1-8                	18317797	        61.4 ns/op
BenchmarkChan1ToN-8                	13170680	        85.3 ns/op
BenchmarkChanNToN-8                	16088103	        72.5 ns/op
BenchmarkChanBufferedRead-8        	49805517	        24.4 ns/op
BenchmarkChanBufferedWrite-8       	46749292	        25.9 ns/op
```
