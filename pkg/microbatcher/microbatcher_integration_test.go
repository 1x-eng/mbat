package microbatcher

import (
	"testing"
	"time"

	"github.com/1x-eng/mbat/pkg/batchprocessor"
	"github.com/1x-eng/mbat/pkg/job"
	"github.com/1x-eng/mbat/pkg/jobqueue"
	"github.com/stretchr/testify/assert"
)

func TestMicroBatcherIntegration(t *testing.T) {
	jobqueue.InitQueue(10)
	jobqueue.InitProcessedJobsCache(5*time.Minute, 10*time.Minute)

	batchProcessor := batchprocessor.NewMockBatchProcessor()

	config := MicroBatcherConfig{
		BatchSize:     1,
		BatchInterval: 1 * time.Millisecond,
	}

	microBatcher := NewMicroBatcher(config, batchProcessor)
	assert.NotNil(t, microBatcher, "NewMicroBatcher should not return nil")

	microBatcher.Start()

	testJob := job.NewJob("test data")
	resultChan := microBatcher.Submit(testJob)
	assert.NotNil(t, resultChan, "Submit should return a result channel")

	result := <-resultChan
	assert.Nil(t, result.Error, "Job should be processed without error")

	cachedJob, found := jobqueue.GetProcessedJobByID(testJob.ID.String())
	assert.True(t, found, "Job should be found in the processed jobs cache")
	assert.NotNil(t, cachedJob, "Cached job should not be nil")

	// Shutdown the MicroBatcher
	microBatcher.Shutdown()
}
