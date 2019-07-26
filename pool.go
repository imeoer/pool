package pool

import (
	"sync/atomic"
	"time"
)

// WorkerPool - use Results to get job result
type WorkerPool struct {
	worker  uint
	remain  uint64
	queue   chan WorkerPoolJob
	Results chan WorkerPoolJob
}

// WorkerPoolJob - need to implement the Do() method
type WorkerPoolJob interface {
	Do()
}

// NewWorkerPool - create workers pool and do job
func NewWorkerPool(jobs []WorkerPoolJob, worker uint, duration time.Duration) *WorkerPool {
	count := len(jobs)
	pool := &WorkerPool{
		worker:  worker,
		remain:  uint64(count),
		queue:   make(chan WorkerPoolJob, count),
		Results: make(chan WorkerPoolJob, count),
	}
	if count <= 0 {
		close(pool.queue)
		close(pool.Results)
		return pool
	}
	for _, job := range jobs {
		pool.queue <- job
	}
	for count := uint(0); count < pool.worker; count++ {
		go func() {
			for {
				job, ok := <-pool.queue
				if !ok {
					return
				}
				job.Do()
				pool.Results <- job
				if atomic.AddUint64(&pool.remain, ^uint64(0)) == 0 {
					close(pool.queue)
					close(pool.Results)
				}
				if duration > 0 {
					time.Sleep(duration)
				}
			}
		}()
	}
	return pool
}
