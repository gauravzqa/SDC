package queue

import (
	"sync"
)

type BlockingBoundedQueue struct {
	queue    chan interface{}
	mutex    sync.Mutex
	notEmpty *sync.Cond
	notFull  *sync.Cond
}

func NewBlockingBoundedQueue(capacity int) *BlockingBoundedQueue {
	bbq := &BlockingBoundedQueue{
		queue: make(chan interface{}, capacity),
	}
	bbq.notEmpty = sync.NewCond(&bbq.mutex)
	bbq.notFull = sync.NewCond(&bbq.mutex)
	return bbq
}

func (bbq *BlockingBoundedQueue) Enqueue(item interface{}) {
	bbq.mutex.Lock()
	defer bbq.mutex.Unlock()

	for len(bbq.queue) == cap(bbq.queue) {
		bbq.notFull.Wait()
	}

	bbq.queue <- item

	bbq.notEmpty.Signal()
}

func (bbq *BlockingBoundedQueue) Dequeue() interface{} {
	bbq.mutex.Lock()
	defer bbq.mutex.Unlock()

	for len(bbq.queue) == 0 {
		bbq.notEmpty.Wait()
	}

	item := <-bbq.queue

	bbq.notFull.Signal()
	return item
}
