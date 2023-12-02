package batchprocessor

import (
	"github.com/1x-eng/mbat/pkg/job"
)

// A contract for BatchProcessor that we DI.
// Im opting for a pointer here, to avoid copying the job object.
type BatchProcessor interface {
	ProcessBatch(jobs []*job.Job) []job.JobResult
}
