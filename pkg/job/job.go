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
