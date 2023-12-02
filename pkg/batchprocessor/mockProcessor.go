package batchprocessor

import (
	"fmt"

	"github.com/1x-eng/mbat/pkg/job"
)

// A mock implementation - to complete main.go
type MockBatchProcessor struct{}

func (m *MockBatchProcessor) ProcessBatch(jobs []*job.Job) []job.JobResult {
	results := make([]job.JobResult, len(jobs))
	for i, j := range jobs {
		// Dummy processing logic
		results[i] = job.JobResult{
			Data:  fmt.Sprintf("Processed job with data: %v", j.Data),
			Error: nil,
		}
	}
	return results
}

func NewMockBatchProcessor() BatchProcessor {
	return &MockBatchProcessor{}
}
