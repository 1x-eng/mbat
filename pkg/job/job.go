package job

type Job struct {
	Data       interface{}
	resultChan chan JobResult
}

func NewJob(data interface{}) *Job {
	return &Job{
		Data:       data,
		resultChan: make(chan JobResult, 1),
	}
}

type JobResult struct {
	Data  interface{}
	Error error
}
