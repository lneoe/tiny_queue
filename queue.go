// queue.go

package main

import (
	"sync"
)

type Message interface {
	GetID() int64
}

type Queue struct {
	closed bool
	queue  []Message

	lock sync.Mutex

	monitor *sync.Cond
}

func NewQ() *Queue {
	q := Queue{queue: nil}
	q.monitor = sync.NewCond(&sync.Mutex{})

	return &q
}

func (q *Queue) Append(msg Message) bool {
	if q.closed {
		return false
	}

	q.lock.Lock()
	defer q.lock.Unlock()

	q.queue = append(q.queue, msg)

	q.monitor.Signal()
	return true
}

func (q *Queue) Close() {
	if q.closed {
		return
	}

	q.closed = true
	q.monitor.Signal()
}

func (q *Queue) Pop() (Message, bool) {
	q.lock.Lock()
	for len(q.queue) > 0 {

		msg := q.queue[0]
		q.queue = q.queue[1:]

		q.lock.Unlock()
		return msg, true
	}
	q.lock.Unlock()

	q.monitor.L.Lock()
	q.monitor.Wait()
	q.monitor.L.Unlock()
	if q.closed {
		return nil, false
	}

	q.lock.Lock()
	msg := q.queue[0]
	q.queue = q.queue[1:]
	q.lock.Unlock()
	return msg, true
}
