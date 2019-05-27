package pool

import (
	"fmt"
	"testing"
	"time"
)

type TestJob struct {
	id int
}

func (job *TestJob) Do() {
	<-time.After(time.Second * 1)
}

func TestRun(t *testing.T) {
	var jobs []WorkerPoolJob
	for i := 1; i <= 100; i++ {
		job := &TestJob{id: i}
		jobs = append(jobs, job)
	}
	pool := NewWorkerPool(jobs, 20, time.Millisecond*10)
	for job := range pool.Results {
		fmt.Printf("Done job %d\n", job.(*TestJob).id)
	}
}
