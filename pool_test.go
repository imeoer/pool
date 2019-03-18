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
	po := New(3)
	for i := 1; i <= 10; i++ {
		job := &TestJob{id: i}
		go po.Put(job)
	}
	for job := range po.Output {
		if po.Idle() {
			break
		}
		fmt.Printf("Done job %d\n", job.(*TestJob).id)
	}
}
