package util

type empty struct{}

type Semaphore chan empty

func NewSemaphore(permits int) Semaphore {
	return make(Semaphore, permits)
}

func (semaphore Semaphore) Acquire() {
	//noinspection GoAssignmentToReceiver
	semaphore <- empty{}
}

func (semaphore Semaphore) Release() {
	<-semaphore
}
