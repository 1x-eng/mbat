package job

import (
	"sync"

	"github.com/google/uuid"
)

type Job struct {
	ID         uuid.UUID
	Data       interface{}
	result     JobResult
	resultChan chan JobResult
	resultOnce sync.Once
}

type JobResult struct {
	Data  interface{}
	Error error
}

func (j *Job) ResultChan() chan JobResult {
	return j.resultChan
}

func (j *Job) SetResult(result JobResult) {
	j.resultOnce.Do(func() {
		j.result = result
		close(j.resultChan)
	})
}

func (j *Job) GetResult() JobResult {
	return j.result
}

func NewJob(data interface{}) *Job {
	return &Job{
		ID:         uuid.New(),
		Data:       data,
		resultChan: make(chan JobResult, 1),
	}
}

// a slice as a temporary storage for all processed jobs.
var processedJobs []*Job

// and a poor man's DLQ. Something quick and not so pretty.
var deadJobs []*Job

func GetProcessedJobs() []*Job {
	return processedJobs
}

func GetProcessedJobByID(jobID uuid.UUID) (*Job, bool) {
	for _, j := range processedJobs {
		if j.ID == jobID {
			return j, true
		}
	}
	return nil, false
}

func StoreProcessedJob(j *Job) {
	processedJobs = append(processedJobs, j)
}

func GetDeadJobs() []*Job {
	return deadJobs
}

func StoreDeadJob(j *Job) {
	deadJobs = append(deadJobs, j)
}
