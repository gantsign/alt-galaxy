package util

type CompletionLatch interface {
	Success()
	Failure()
	Await() bool
}

type completionLatchImpl struct {
	count   int
	channel chan bool
}

func NewCompletionLatch(count int) CompletionLatch {
	return completionLatchImpl{count, make(chan bool, count)}
}

func (latch completionLatchImpl) Success() {
	latch.channel <- true
}

func (latch completionLatchImpl) Failure() {
	latch.channel <- false
}

func (latch completionLatchImpl) Await() bool {
	result := true
	for i := 0; i < latch.count; i++ {
		result = result && <-latch.channel
	}
	return result
}
