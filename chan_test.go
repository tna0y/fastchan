package fastchan

import (
	"runtime"
	"sync"
	"testing"
)

func BenchmarkChan1To1(b *testing.B) {
	ch := make(chan int, myMax(uint64(b.N), 2))
	b.ResetTimer()
	go func() {
		for i := 0; i < b.N; i++ {
			ch <- 1
		}
	}()
	for i := 0; i < b.N; i++ {
		<-ch
	}
}

func BenchmarkChanNTo1(b *testing.B) {
	ch := make(chan int, myMax(uint64(b.N), 2))
	b.ResetTimer()
	cores := runtime.NumCPU()
	perGoro := b.N / cores
	for i := 0; i < cores; i++ {
		go func() {
			for j := 0; j < perGoro; j++ {
				ch <- 1
			}
		}()
	}
	for i := 0; i < perGoro*cores; i++ {
		<-ch
	}
}

func BenchmarkChan1ToN(b *testing.B) {
	ch := make(chan int, myMax(uint64(b.N), 2))
	wg := sync.WaitGroup{}
	cores := runtime.NumCPU()
	perGoro := b.N / cores
	b.ResetTimer()
	for i := 0; i < cores; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < perGoro; j++ {
				<-ch
			}
			wg.Done()
		}()
	}
	for i := 0; i < perGoro*cores; i++ {
		ch <- i
	}
	wg.Wait()
}

func BenchmarkChanNToN(b *testing.B) {
	ch := make(chan int, myMax(uint64(b.N), 2))
	wg := sync.WaitGroup{}
	cores := runtime.NumCPU()
	perGoro := b.N / cores
	b.ResetTimer()

	for i := 0; i < cores; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < perGoro; j++ {
				ch <- j
			}
		}()
		go func() {
			for j := 0; j < perGoro; j++ {
				<-ch
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
