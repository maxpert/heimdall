package heimdall


type WorkerPool interface {
	Schedule(task func())
}

type channelWorkerPool struct {
	worker    chan func()
	semaphore chan struct{}
}

var semaphorePayload = struct{}{}

func (c *channelWorkerPool) Schedule(task func()) {
	select {
	case c.semaphore <- semaphorePayload:
		go c.runWorkers(task)
	case c.worker <- task:
	}
}

func (c *channelWorkerPool) runWorkers(task func()) {
	defer (func() { <-c.semaphore })()
	for {
		task()
		task = <-c.worker
	}
}

func NewWorkerPool(size int) WorkerPool {
	return &channelWorkerPool{
		worker:    make(chan func()),
		semaphore: make(chan struct{}, size),
	}
}
