package atomic_go //import "github.com/Cybergenik/atomic-go"

import (
	"sync"
	"sync/atomic"
)

type AtomicQueue[T any] struct {
	head     uint32
	tail     uint32
	size     uint32
	len      uint32
	headLock sync.Mutex
	tailLock sync.Mutex
	buffer   []*T
}

func (q *AtomicQueue[T]) Init(size uint32) {
	q.headLock = sync.Mutex{}
	q.tailLock = sync.Mutex{}
	q.size = size
	q.len = 0
	q.head = 0
	q.tail = 0
	q.buffer = make([]*T, size)
}

func (q *AtomicQueue[T]) Push(item *T) bool {
	q.tailLock.Lock()
	defer q.tailLock.Unlock()
	if q.tail == q.size {
		q.tail %= q.size
	}
	if q.buffer[q.tail] != nil {
		return false
	}
	q.buffer[q.tail] = item
	atomic.StoreUint32(&q.len, q.len+1)
	q.tail++
	return true
}

func (q *AtomicQueue[T]) Pop() *T {
	q.headLock.Lock()
	defer q.headLock.Unlock()
	if q.head == q.size {
		q.head %= q.size
	}
	t := q.buffer[q.head]
	if t == nil {
		return nil
	}
	q.buffer[q.head] = nil
	atomic.StoreUint32(&q.len, q.len-1)
	q.head++
	return t
}

func (q *AtomicQueue[T]) Len() uint32 {
	return atomic.LoadUint32(&q.len)
}

func (q *AtomicQueue[T]) Size() uint32 {
	return q.size
}
