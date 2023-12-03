package jobqueue

import (
	"log"
	"sync"

	"github.com/1x-eng/mbat/pkg/job"
	"github.com/google/uuid"
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
	log.Printf("Enqueued job %s\n", j.ID)
}

func Dequeue(jobID uuid.UUID) {
	mapLock.Lock()
	defer mapLock.Unlock()
	delete(jobMap, jobID.String())
	log.Printf("Dequeued job %s\n", jobID)
}

func FindJobByID(jobID string) (*job.Job, bool) {
	mapLock.RLock()
	j, exists := jobMap[jobID]
	mapLock.RUnlock()
	return j, exists
}

func GetQueuedJobIds() []string {
	mapLock.RLock()
	defer mapLock.RUnlock()

	var jobIDs []string
	for id := range jobMap {
		jobIDs = append(jobIDs, id)
	}
	return jobIDs
}

func DrainQueue() {
	// Ideally, this is where I want to use a DLQ, but short on time for now. Will revisit.
	for len(JobQueue) > 0 {
		<-JobQueue
	}
	close(JobQueue)
	log.Println("Job queue drained & closed")
}
