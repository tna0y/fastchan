package fastchan

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

func TestMSQueueBasic(t *testing.T) {
	n := 1000
	q := newMSQueue()
	for i := 0; i < n; i++ {
		q.push(i)
	}
	for i := 0; i < n; i++ {
		if v, _ := q.pop(); i != v {
			t.Fatal("Pop not equal")
		}
	}
}

func TestMSQueueConcurrent(t *testing.T) {
	n := 100000
	q := newMSQueue()
	go func() {
		for i := 0; i < n; i++ {
			q.push(i)
		}
	}()
	for i := 0; i < n; i++ {
		if v, _ := q.pop(); i != v {
			t.Fatal("Pop not equal")
		}
	}
}

func TestMSQueuePreserveOrder(t *testing.T) {
	n := 10000000
	nGoros := runtime.NumCPU()
	if nGoros > 10 {
		nGoros = 10
	}

	q := newMSQueue()

	start := sync.WaitGroup{}
	stop := sync.WaitGroup{}
	start.Add(1)
	stop.Add(nGoros)
	var done uint32

	for g := 0; g < nGoros; g++ {
		go func(g int) {
			start.Wait()
			for i := 0; i < n; i += 10 {
				q.push(i + g)
			}
			atomic.AddUint32(&done, 1)
		}(g)

		go func(g int) {
			start.Wait()
			lasts := make(map[int]int)
			for {

				item, ok := q.pop()
				if !ok {
					if atomic.LoadUint32(&done) == uint32(nGoros) {
						break
					}
					runtime.Gosched()
					continue
				}
				rem := item % 10
				if v, ok := lasts[rem]; !ok || v < item {
					lasts[rem] = item
				} else {
					t.Fatal("Fail", ok, v, item)
				}

			}
			stop.Done()
		}(g)
	}
	start.Done()
	stop.Wait()

}
