package jobqueue

import (
	"log"

	"github.com/1x-eng/mbat/pkg/job"
	"github.com/google/uuid"
)

var JobQueue chan *job.Job

func InitQueue(queueSize int) {
	JobQueue = make(chan *job.Job, queueSize)
}

func Enqueue(job *job.Job) {
	JobQueue <- job
}

func Dequeue() *job.Job {
	return <-JobQueue
}

func FindJobByID(jobID uuid.UUID) (*job.Job, bool) {
	for j := range JobQueue {
		if j.ID == jobID {
			return j, true
		}
	}

	return nil, false
}

func DrainQueue() {
	for len(JobQueue) > 0 {
		job.StoreDeadJob(<-JobQueue)
	}
	close(JobQueue)
	log.Printf("Job queue drained, any remaining jobs are hydrated into the dead job queue. Job queue closed.")
}
