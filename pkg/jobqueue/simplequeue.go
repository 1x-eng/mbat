package jobqueue

import (
	"log"
	"sync"

	"github.com/1x-eng/mbat/pkg/job"
)

var (
	JobQueue chan *job.Job
	jobMap   map[string]*job.Job
	mapLock  sync.RWMutex
)

func InitQueue(queueSize int) {
	JobQueue = make(chan *job.Job, queueSize)
	jobMap = make(map[string]*job.Job)
}

func Enqueue(j *job.Job) {
	mapLock.Lock()
	jobMap[j.ID.String()] = j
	mapLock.Unlock()

	JobQueue <- j
}

func Dequeue() *job.Job {
	j := <-JobQueue
	mapLock.Lock()
	delete(jobMap, j.ID.String())
	mapLock.Unlock()
	return j
}

func FindJobByID(jobID string) (*job.Job, bool) {
	mapLock.RLock()
	j, exists := jobMap[jobID]
	mapLock.RUnlock()
	return j, exists
}

func DrainQueue() {
	// Ideally, this is where I want to use a DLQ, but short on time for now. Will revisit.
	for len(JobQueue) > 0 {
		<-JobQueue
	}
	close(JobQueue)
	log.Println("Job queue drained & closed")
}
