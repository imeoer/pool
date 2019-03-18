package pool

// Pool - pool object
type Pool struct {
	size    uint
	workers chan Job
	remain  uint
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
					pool.remain--
				}()
			}
		}()
	}
	return pool
}

// Put - put the job to workers pool
func (pool *Pool) Put(job Job) {
	pool.workers <- job
	pool.remain++
}

// Idle - check workers pool idle status
func (pool *Pool) Idle() bool {
	return pool.remain == 0
}
