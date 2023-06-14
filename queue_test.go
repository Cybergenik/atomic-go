package atomic_go //import "github.com/Cybergenik/atomic-go"

import (
	"sync"
	"testing"
)

func TestQueueOrder(t *testing.T) {
	q := &AtomicQueue[string]{}
	q.Init(4)
	if q.Size() != 4 {
		t.Errorf("Incorrect Queue Size: %+v", q)
	}
	strs := []string{"Together", "We", "Must", "Bus"}
	q.Push(&strs[0])
	q.Push(&strs[1])
	q.Push(&strs[2])
	q.Push(&strs[3])
	for _, s := range strs {
		qString := q.Pop()
		if s != *qString {
			t.Errorf("Queue Invalid Order: %s != %s", s, *qString)
		}
	}
}

func TestQueueOversized(t *testing.T) {
	q := &AtomicQueue[string]{}
	q.Init(10)
	if q.Size() != 10 {
		t.Errorf("Incorrect Queue Size: %+v", q)
	}
	strs := []string{"Together", "We", "Must", "Bus"}
	q.Push(&strs[0])
	q.Push(&strs[1])
	q.Push(&strs[2])
	q.Push(&strs[3])
	if q.Len() != 4 {
		t.Errorf("Incorrect number of items in Queue: %+v", q)
	}
	for _, s := range strs {
		qString := q.Pop()
		if s != *qString {
			t.Errorf("Queue Invalid Order: %s != %s", s, *qString)
		}
	}
	if q.Len() != 0 {
		t.Errorf("Incorrect number of items in Queue: %+v", q)
	}
}

func TestQueueLarge(t *testing.T) {
	q := &AtomicQueue[string]{}
	q.Init(9)
	if q.Size() != 9 {
		t.Errorf("Incorrect Queue Size: %+v", q)
	}
	strs := []string{"Together", "We", "Must", "Bus", "And", "Hopefully", "Vape", "A", "Bit"}
	q.Push(&strs[0])
	q.Push(&strs[1])
	q.Push(&strs[2])
	q.Push(&strs[3])
	q.Push(&strs[4])
	q.Push(&strs[5])
	q.Push(&strs[6])
	q.Push(&strs[7])
	q.Push(&strs[8])
	if q.Len() != 9 {
		t.Errorf("Incorrect Queue Size: %+v", q)
	}
	for _, s := range strs {
		qString := q.Pop()
		if s != *qString {
			t.Errorf("Queue Invalid Order: %s != %s", s, *qString)
		}
	}
	if q.Len() != 0 {
		t.Errorf("Incorrect Queue Size: %+v", q)
	}
}

func TestQueueTooManyItems(t *testing.T) {
	q := &AtomicQueue[string]{}
	q.Init(9)
	strs := []string{"Together", "We", "Must", "Bus", "And", "Hopefully", "Vape", "A", "Bit"}
	q.Push(&strs[0])
	q.Push(&strs[1])
	q.Push(&strs[2])
	q.Push(&strs[3])
	q.Push(&strs[4])
	q.Push(&strs[5])
	q.Push(&strs[6])
	q.Push(&strs[7])
	q.Push(&strs[8])
	if q.Len() != 9 {
		t.Errorf("Incorrect number of items in Queue: %+v", q)
	}
	if q.Push(&strs[0]) {
		t.Errorf("Pushed too Many Items to Queue: %+v", q)
	}
	for _, s := range strs {
		qString := q.Pop()
		if s != *qString {
			t.Errorf("Queue Invalid Order: %s != %s", s, *qString)
		}
	}
	if q.Len() != 0 {
		t.Errorf("Incorrect Queue Size: %+v", q)
	}
}

func TestQueueNil(t *testing.T) {
	q := &AtomicQueue[string]{}
	q.Init(2)
	strs := []string{"Hello", "World"}
	q.Push(&strs[0])
	q.Push(&strs[1])
	if q.Len() != 2 {
		t.Errorf("Incorrect number of items in Queue: %+v", q)
	}
	q.Pop()
	q.Pop()
	if q.Len() != 0 {
		t.Errorf("Incorrect number of items in Queue: %+v", q)
	}
	qString := q.Pop()
	if qString != nil {
		t.Errorf("Received Non Nil Value: %s", *qString)
	}
}

func TestQueuePopAdd(t *testing.T) {
	q := &AtomicQueue[string]{}
	q.Init(2)
	strs := []string{"Hello", "World", "Coom", "Zoom"}
	qString := ""
	q.Push(&strs[0])
	qString = *q.Pop()
	if qString != strs[0] {
		t.Errorf("Queue Invalid Order: %s != %s", strs[0], qString)
	}
	q.Push(&strs[1])
	qString = *q.Pop()
	if qString != strs[1] {
		t.Errorf("Queue Invalid Order: %s != %s", strs[1], qString)
	}
	q.Push(&strs[2])
	qString = *q.Pop()
	if qString != strs[2] {
		t.Errorf("Queue Invalid Order: %s != %s", strs[2], qString)
	}
	q.Push(&strs[3])
	qString = *q.Pop()
	if qString != strs[3] {
		t.Errorf("Queue Invalid Order: %s != %s", strs[3], qString)
	}
}

func TestQueueParallelSmall(t *testing.T) {
	q := &AtomicQueue[int]{}
	N := 100
	q.Init(uint32(N))
	wg := sync.WaitGroup{}
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(i int) {
			q.Push(&i)
			wg.Done()
		}(i)
	}
	wg.Wait()
	seen := map[int]bool{}
	if q.Len() != uint32(N) {
		t.Errorf("Queue Size Wrong: %+v", q)
	}
	for t := q.Pop(); t != nil; t = q.Pop() {
		seen[*t] = true
	}
	if len(seen) != N {
		t.Errorf("Pop Count Wrong: %+v", q)
	}
}

func TestQueueParallelLarge(t *testing.T) {
	q := &AtomicQueue[int]{}
	N := 10_000
	q.Init(uint32(N))
	wg := sync.WaitGroup{}
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(i int) {
			q.Push(&i)
			wg.Done()
		}(i)
	}
	wg.Wait()
	seen := map[int]bool{}
	if q.Len() != uint32(N) {
		t.Errorf("Queue Size Wrong: %+v", q)
	}
	for t := q.Pop(); t != nil; t = q.Pop() {
		seen[*t] = true
	}
	if len(seen) != N {
		t.Errorf("Pop Count Wrong: %+v", q)
	}
}

// BENCHMARKS:
//
// Comparing my Atomic Queue vs buffered channel
// 
///////
func runQueuePush(N int, q *AtomicQueue[int], wg *sync.WaitGroup) {
	for i := 0; i < N; i++ {
		go func(i int) {
			q.Push(&i)
			wg.Done()
		}(i)
	}
}

func runQueuePop(N int, q *AtomicQueue[int], wg *sync.WaitGroup) {
	for i := 0; i < N; i++ {
		go func(i int) {
			q.Pop()
			wg.Done()
		}(i)
	}
}

func BenchmarkQueue(b *testing.B) {
    for i:=0; i < b.N; i++ {
        N := 10_000
        q := &AtomicQueue[int]{}
        q.Init(uint32(N))
        wg := &sync.WaitGroup{}
		wg.Add(N/2)
        runQueuePush(N/2, q, wg)
		wg.Add(N/2)
        runQueuePop(N/2, q, wg)
		wg.Add(N/2)
        runQueuePush(N/2, q, wg)
        wg.Wait()
    }
}

func runChanPush(N int, c chan *int, wg *sync.WaitGroup) {
	for i := 0; i < N; i++ {
		go func(i int) {
            c <- &i
			wg.Done()
		}(i)
	}
}

func runChanPop(N int, c chan *int, wg *sync.WaitGroup) {
	for i := 0; i < N; i++ {
		go func(i int) {
            <-c
			wg.Done()
		}(i)
	}
}

func BenchmarkChan(b *testing.B) {
    for i:=0; i < b.N; i++ {
        N := 10_000
        c := make(chan *int, N)
        wg := &sync.WaitGroup{}
		wg.Add(N/2)
        runChanPush(N/2, c, wg)
		wg.Add(N/2)
        runChanPop(N/2, c, wg)
		wg.Add(N/2)
        runChanPush(N/2, c, wg)
        wg.Wait()
    }
}
