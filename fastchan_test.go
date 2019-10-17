package fastchan

import (
	"runtime"
	"testing"
)

func TestIntRing(t *testing.T) {
	n := 1000
	var rb *FastChan
	rb = New(uint64(2))
	go func() {
		for i := 0; i < n; i++ {
			err := rb.Put(i)
			if err != nil {
				t.Fatal(err)
			}
		}
	}()
	for i := 0; i < n; i++ {
		v, err := rb.Get()
		if err != nil {
			t.Fatal(err)
		}
		if v != i {
			t.Fatal("Ti pidor")
		}
	}

}

func myMax(a, b uint64) uint64 {
	if a < b {
		return b
	}
	return a
}

func BenchmarkIntRingBuf1To1(b *testing.B) {
	ch := New(myMax(uint64(b.N), 2))
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

func BenchmarkIntRingBufNTo1(b *testing.B) {
	ch := New(myMax(uint64(b.N), 2))
	b.ResetTimer()
	cores := runtime.NumCPU()
	perGoro := b.N / cores
	for i := 0; i < cores; i++ {
		go func() {
			for j := 0; j < b.N; j++ {
				ch.Put(1)
			}
		}()
	}
	for i := 0; i < perGoro*cores; i++ {
		ch.Get()
	}
}
