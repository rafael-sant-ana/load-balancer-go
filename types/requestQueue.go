package types

import (
	"fmt"
	"sync"
)

type RequestQueue struct {
	Requests   []*RequestEvent
	TopChanged chan *RequestEvent
	mu         sync.Mutex
}

// Constructor (boa prática)

func NewRequestQueue() *RequestQueue {
	return &RequestQueue{
		Requests:   make([]*RequestEvent, 0),
		TopChanged: make(chan *RequestEvent, 1), // buffer 1
	}
}

// Enqueue adds an element to the end of the queue.
func (q *RequestQueue) Enqueue(value *RequestEvent) {
	q.mu.Lock()
	q.Requests = append(q.Requests, value)
	shouldUpdate := len(q.Requests) == 1
	top := q.Requests[0]
	q.mu.Unlock()

	if shouldUpdate {
		select {
		case q.TopChanged <- top:
		default:
		}
	}
}

// Dequeue removes and returns the first element of the queue.
func (q *RequestQueue) Dequeue() (*RequestEvent, error) {
	q.mu.Lock()
	if len(q.Requests) == 0 {
		return nil, fmt.Errorf("Queue is empty")
	}
	oldTop := q.Requests[0]

	value := (q.Requests)[0]
	q.Requests[0] = nil
	q.Requests = q.Requests[1:]

	var newTop *RequestEvent = nil
	if len(q.Requests) > 0 {
		newTop = q.Requests[0]
	}
	shouldUpdate := oldTop != newTop

	q.mu.Unlock()
	if shouldUpdate {
		select {
		case q.TopChanged <- newTop:
		default:
		}
	}
	return value, nil
}

func (q *RequestQueue) Top() *RequestEvent {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.Requests) == 0 {
		return nil
	}
	return (q.Requests)[0]
}
