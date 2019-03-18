package pool

import (
	"sync/atomic"
)

// Pool - pool object
type Pool struct {
	size    uint
	remain  uint64
	workers chan Job
	Output  chan Job
}

// Job - need to implement the Do() method
type Job interface {
	Do()
}

// New - create workers pool
func New(size uint) *Pool {
	pool := &Pool{
		size:    size,
		workers: make(chan Job, size),
		Output:  make(chan Job),
	}
	for count := uint(0); count < pool.size; count++ {
		go func() {
			for {
				job := <-pool.workers
				job.Do()
				go func() {
					pool.Output <- job
					atomic.AddUint64(&pool.remain, ^uint64(0))
				}()
			}
		}()
	}
	return pool
}

// Put - put the job to workers pool
func (pool *Pool) Put(job Job) {
	pool.workers <- job
	atomic.AddUint64(&pool.remain, 1)
}

// Idle - check workers pool idle status
func (pool *Pool) Idle() bool {
	remain := atomic.LoadUint64(&pool.remain)
	return remain == 0
}
