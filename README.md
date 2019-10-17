# Fastchan
Go channels, but faster
## When to use
## Limitations
* Only buffered channels are supported (for now)
* Minimum buffer size is 2
* Buffer size is rounded up to the nearest power of two
## How to use
## Benchmarks
very fast
## Credits
* Implementation idea [1024cores.net](http://www.1024cores.net/home/lock-free-algorithms/queues/bounded-mpmc-queue)
* Original implementation [Workiva/go-datastructures](https://github.com/Workiva/go-datastructures/blob/master/queue/ring.go)