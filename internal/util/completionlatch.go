package util

import (
	"sync"
	"sync/atomic"
)

type CompletionLatch interface {
	TaskAdded()
	Success()
	Failure()
	Await() bool
}

type completionLatchImpl struct {
	remaining      *int32
	result         *int32
	updateLock     *sync.Mutex
	completionLock *sync.Mutex
}

func NewCompletionLatch(count int) CompletionLatch {
	count32 := int32(count)
	result32 := int32(1)

	latch := completionLatchImpl{
		remaining:      &count32,
		result:         &result32,
		updateLock:     &sync.Mutex{},
		completionLock: &sync.Mutex{},
	}
	latch.completionLock.Lock()
	return latch
}

func (latch completionLatchImpl) TaskAdded() {
	latch.updateLock.Lock()

	atomic.AddInt32(latch.remaining, 1)

	latch.updateLock.Unlock()
}

func (latch completionLatchImpl) Success() {
	latch.updateLock.Lock()

	atomic.AddInt32(latch.remaining, -1)

	if atomic.LoadInt32(latch.remaining) == 0 {
		latch.completionLock.Unlock()
	}

	latch.updateLock.Unlock()
}

func (latch completionLatchImpl) Failure() {
	latch.updateLock.Lock()

	atomic.AddInt32(latch.remaining, -1)
	atomic.StoreInt32(latch.result, 0)

	if atomic.LoadInt32(latch.remaining) == 0 {
		latch.completionLock.Unlock()
	}

	latch.updateLock.Unlock()
}

func (latch completionLatchImpl) Await() bool {
	latch.completionLock.Lock()

	result := atomic.LoadInt32(latch.result) == 1

	latch.completionLock.Unlock()

	return result
}
