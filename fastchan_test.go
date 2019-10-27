package fastchan

import (
	"runtime"
	"sync"
	"testing"
)

//
// Tests
//

func TestBasic(t *testing.T) {
	n := 10000
	var rb *FastChan
	rb = New(uint32(1))
	go func() {
		for i := 0; i < n; i++ {
			rb.Put(i)
		}
	}()
	for i := 0; i < n; i++ {
		v := rb.Get()
		if v != i {
			t.Fatal("fail")
		}
	}

}

func TestBufferSizeOne(t *testing.T) {
	n := 1000
	var rb *FastChan
	rb = New(uint32(1))
	go func() {
		for i := 0; i < n; i++ {
			rb.Put(i)
		}
	}()
	for i := 0; i < n; i++ {
		v := rb.Get()
		if v != i {
			t.Fatal("fail")
		}
	}
}

//
// Benchmarks
//

func myMax(a, b uint32) uint32 {
	if a < b {
		return b
	}
	return a
}

func BenchmarkFastChanBufferedRead(b *testing.B) {
	ch := New(uint32(b.N))
	for i := 0; i < b.N; i++ {
		ch.Put(1)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ch.Get()
	}
}

func BenchmarkFastChanBufferedWrite(b *testing.B) {
	ch := New(uint32(b.N))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ch.Put(1)
	}
}

func BenchmarkFastChan1To1(b *testing.B) {
	ch := New(myMax(uint32(b.N), 2))
	b.ResetTimer()
	go func() {
		for i := 0; i < b.N; i++ {
			ch.Put(1)
		}
	}()
	for i := 0; i < b.N; i++ {
		ch.Get()
	}
}

func BenchmarkFastChanNTo1(b *testing.B) {
	ch := New(myMax(uint32(b.N), 2))
	cores := runtime.NumCPU()
	perGoro := b.N / cores
	b.ResetTimer()
	for i := 0; i < cores; i++ {
		go func() {
			for j := 0; j < perGoro; j++ {
				ch.Put(1)
			}
		}()
	}
	for i := 0; i < perGoro*cores; i++ {
		ch.Get()
	}
}

func BenchmarkFastChan1ToN(b *testing.B) {
	ch := New(myMax(uint32(b.N), 2))
	wg := sync.WaitGroup{}
	cores := runtime.NumCPU()
	perGoro := b.N / cores
	b.ResetTimer()
	for i := 0; i < cores; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < perGoro; j++ {
				ch.Get()
			}
			wg.Done()
		}()
	}
	for i := 0; i < perGoro*cores; i++ {
		ch.Put(i)
	}
	wg.Wait()
}

func BenchmarkFastChanNToN(b *testing.B) {
	ch := New(myMax(uint32(b.N), 2))
	wg := sync.WaitGroup{}
	cores := runtime.NumCPU()
	perGoro := b.N / cores
	b.ResetTimer()

	for i := 0; i < cores; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < perGoro; j++ {
				ch.Put(j)
			}
		}()
		go func() {
			for j := 0; j < perGoro; j++ {
				ch.Get()
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
