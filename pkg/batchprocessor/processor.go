package batchprocessor

import "github.com/1x-eng/mbat/pkg/job"

// A contract for BatchProcessor that we DI.
type BatchProcessor interface {
	ProcessBatch(jobs []job.Job) []job.JobResult
}
