package job

type Job struct {
	Data       interface{}
	resultChan chan JobResult
}

type JobResult struct {
	Data  interface{}
	Error error
}

func (j *Job) ResultChan() chan JobResult {
	return j.resultChan
}

func (j *Job) SetResult(result JobResult) {
	j.resultChan <- result
	close(j.resultChan)
}

func NewJob(data interface{}) *Job {
	return &Job{
		Data:       data,
		resultChan: make(chan JobResult, 1),
	}
}
