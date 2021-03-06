package fastchan

import (
	"runtime"
	"sync/atomic"
	"unsafe"
)

func roundUp(v uint32) uint32 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v |= v >> 32
	v++
	return v
}

type node struct {
	position uint64
	data     interface{}
}

const (
	nodePtrSize        = unsafe.Sizeof(&node{})
	maskHigh    uint64 = 1 << 63
	maskLow     uint64 = (1 << 63) - 1
)

type FastChan struct {
	_padding0 [8]uint64
	queue     uint64
	_padding1 [8]uint64
	dequeue   uint64
	_padding2 [8]uint64
	mask      uint64
	closed    uint32
	_padding3 [8]uint64
	nodes     []*node
	nodePtr   uintptr
	_padding4 [8]uint64
}

func (fc *FastChan) init(size uint32) {
	size = roundUp(size)
	fc.nodes = make([]*node, size)
	fc.nodePtr = uintptr(unsafe.Pointer(&fc.nodes[0]))

	for i := uint32(0); i < size; i++ {
		fc.nodes[i] = &node{position: uint64(i)}
	}
	fc.mask = uint64(size - 1) // so we don't have to do this with every put/get operation
}

// Put adds the provided item to the queue.  If the queue is full, this
// call will block until an item is added to the queue or Close is called
// on the queue.  An error will be returned if the queue is disposed.
func (fc *FastChan) Put(item interface{}) {
	fc.put(item, false)
}

// Offer adds the provided item to the queue if there is space.  If the queue
// is full, this call will return false.  An error will be returned if the
// queue is disposed.
func (fc *FastChan) TryPut(item interface{}) bool {
	return fc.put(item, true)
}

func (fc *FastChan) put(item interface{}, offer bool) bool {

	var (
		n       *node
		pos     uint64
		cPos    uint64
		full    bool
		itemPos uint64
	)
	for {
		if atomic.LoadUint32(&fc.closed) == 1 {
			panic("Put on closed fastchan")
		}
		pos = fc.queue

		// The same as  n = fc.nodes[pos&fc.mask] but without bounds check
		n = *(**node)(unsafe.Pointer(fc.nodePtr + uintptr(pos&fc.mask)*nodePtrSize))

		cPos = n.position
		full = cPos&maskHigh == maskHigh
		itemPos = cPos & maskLow
		if !full && itemPos == pos && atomic.CompareAndSwapUint64(&fc.queue, pos, pos+1) {
			break
		}

		if offer {
			return false
		}
		runtime.Gosched()
	}

	n.data = item
	atomic.StoreUint64(&n.position, cPos|maskHigh)
	return true
}

func (fc *FastChan) Get() interface{} {

	var (
		n       *node
		pos     uint64
		cPos    uint64
		full    bool
		itemPos uint64
	)
	for {
		if atomic.LoadUint32(&fc.closed) == 1 {
			panic("Get on closed fastchan")
		}
		pos = fc.dequeue

		// The same as  n = fc.nodes[pos&fc.mask] but without bounds check
		n = *(**node)(unsafe.Pointer(fc.nodePtr + uintptr(pos&fc.mask)*nodePtrSize))

		cPos = n.position
		full = cPos&maskHigh == maskHigh
		itemPos = cPos & maskLow
		if full && itemPos == pos && atomic.CompareAndSwapUint64(&fc.dequeue, pos, pos+1) {
			break
		}

		runtime.Gosched()
	}
	data := n.data
	n.data = nil
	atomic.StoreUint64(&n.position, (pos+fc.mask+1)&maskLow)
	return data
}

// Len returns the number of items in the queue.
func (fc *FastChan) Len() uint32 {
	return uint32(fc.queue - fc.dequeue)
}

// Cap returns the capacity of this ring buffer.
func (fc *FastChan) Cap() uint64 {
	return uint64(len(fc.nodes))
}

// Close will dispose of this queue and free any blocked threads
// in the Put and/or Get methods.  Calling those methods on a disposed
// queue will return an error.
func (fc *FastChan) Close() {
	fc.closed = 1
}

// IsClosed will return a bool indicating if this queue has been
// closed.
func (fc *FastChan) IsClosed() bool {
	return fc.closed == 1
}

func New(size uint32) *FastChan {
	rb := &FastChan{}
	rb.init(size)
	return rb
}

// Michael-Scott unbounded queue
// http://www.cs.rochester.edu/~scott/papers/1996_PODC_queues.pdf?
type msqueue struct {
	head *msnode
	tail *msnode
}

type msnode struct {
	val  int
	next *msnode
}

func newMSQueue() *msqueue {
	item := &msnode{}
	return &msqueue{
		head: item,
		tail: item,
	}
}

func (q *msqueue) push(val int) {
	// TODO: use thread-local struct cache
	newTail := &msnode{val: val}
	for {
		tail := q.tail
		if atomic.CompareAndSwapPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&tail.next)),
			nil,
			unsafe.Pointer(newTail)) {
			atomic.CompareAndSwapPointer(
				(*unsafe.Pointer)(unsafe.Pointer(&q.tail)),
				unsafe.Pointer(tail),
				unsafe.Pointer(newTail))
			return
		} else {
			atomic.CompareAndSwapPointer(
				(*unsafe.Pointer)(unsafe.Pointer(&q.tail)),
				unsafe.Pointer(tail),
				unsafe.Pointer(tail.next))
		}
	}
}

func (q *msqueue) pop() (int, bool) {
	for {
		head := q.head
		tail := q.tail
		nextHead := head.next
		if head == tail {
			if nextHead == nil {
				return 0, false
			} else {
				atomic.CompareAndSwapPointer(
					(*unsafe.Pointer)(unsafe.Pointer(&q.tail)),
					unsafe.Pointer(tail),
					unsafe.Pointer(tail.next))
			}
		} else {
			res := nextHead.val
			if atomic.CompareAndSwapPointer(
				(*unsafe.Pointer)(unsafe.Pointer(&q.head)),
				unsafe.Pointer(head),
				unsafe.Pointer(nextHead)) {
				return res, true
			}
		}
	}
}
