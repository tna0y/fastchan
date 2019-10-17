package fastchan

import (
	"sync/atomic"
)

import "errors"

var (
	// ErrDisposed is returned when an operation is performed on a disposed
	// queue.
	ErrDisposed = errors.New(`queue: disposed`)

	// ErrTimeout is returned when an applicable queue operation times out.
	ErrTimeout = errors.New(`queue: poll timed out`)

	// ErrEmptyQueue is returned when an non-applicable queue operation was called
	// due to the queue's empty item state
	ErrEmptyQueue = errors.New(`queue: empty queue`)
)

// roundUp takes a uint64 greater than 0 and rounds it up to the next
// power of 2.
func roundUp(v uint64) uint64 {
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

type nodes []*node

type FastChan struct {
	_padding0    [8]uint64
	queue        uint64
	_padding1    [8]uint64
	dequeue      uint64
	_padding2    [8]uint64
	mask, closed uint64
	_padding3    [8]uint64
	nodes        nodes
	_padding4    [8]uint64
	semPut       uint32
	_paddng5     [8]uint64
	semGet       uint32
	_padding5    [8]uint64
}

func (fc *FastChan) init(size uint64) {
	size = roundUp(size)
	fc.nodes = make(nodes, size)
	for i := uint64(0); i < size; i++ {
		fc.nodes[i] = &node{position: i}
	}
	fc.semPut = uint32(size)
	fc.mask = size - 1 // so we don't have to do this with every put/get operation
}

// Put adds the provided item to the queue.  If the queue is full, this
// call will block until an item is added to the queue or Close is called
// on the queue.  An error will be returned if the queue is disposed.
func (fc *FastChan) Put(item interface{}) error {
	_, err := fc.put(item, false)
	return err
}

// Offer adds the provided item to the queue if there is space.  If the queue
// is full, this call will return false.  An error will be returned if the
// queue is disposed.
func (fc *FastChan) TryPut(item interface{}) (bool, error) {
	return fc.put(item, true)
}

// We avoid using atomic loads and stores, since they offer no memory barriers. We can avoid calls
// See https://github.com/golang/go/issues/5045
func (fc *FastChan) put(item interface{}, offer bool) (bool, error) {

	var n *node
	var pos uint64
	for {
		if fc.closed == 1 {
			return false, ErrDisposed
		}
		pos = fc.queue
		n = fc.nodes[pos&fc.mask]
		if n.position == pos && atomic.CompareAndSwapUint64(&fc.queue, pos, pos+1) {
			break
		}

		if offer {
			return false, nil
		}

		// We allow semaphores to go over to higher values. This causes some spinning but eventually leads to threads
		// stalling and waiting to
		semacquire(&fc.semPut)

	}

	n.data = item
	n.position = pos + 1
	semrelease(&fc.semGet, false, 0)
	return true, nil
}

func (fc *FastChan) Get() (interface{}, error) {

	var (
		n   *node
		pos uint64
	)
	for {
		if fc.closed == 1 {
			return nil, ErrDisposed
		}
		pos = fc.dequeue
		n = fc.nodes[pos&fc.mask]

		if n.position == pos+1 && atomic.CompareAndSwapUint64(&fc.dequeue, pos, pos+1) {
			break
		}

		semacquire(&fc.semGet)
	}
	data := n.data
	n.data = nil
	n.position = pos + fc.mask + 1
	semrelease(&fc.semPut, false, 0)
	return data, nil
}

// Len returns the number of items in the queue.
func (fc *FastChan) Len() uint64 {
	return atomic.LoadUint64(&fc.queue) - atomic.LoadUint64(&fc.dequeue)
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

func New(size uint64) *FastChan {
	rb := &FastChan{}
	rb.init(size)
	return rb
}
