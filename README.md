# Fastchan
Go channels, but faster
## When to use
## Limitations
* Only buffered channels are supported (for now)
* Buffer size is rounded up to the nearest power of two
## How to use
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