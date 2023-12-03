package microbatcher

import (
	"log"
	"time"

	"github.com/1x-eng/mbat/pkg/batchprocessor"
	"github.com/1x-eng/mbat/pkg/job"
	"github.com/1x-eng/mbat/pkg/jobqueue"
)

type MicroBatcherConfig struct {
	BatchSize     int
	BatchInterval time.Duration
}

type MicroBatcher struct {
	scheduler *BatchScheduler
}

func NewMicroBatcher(config MicroBatcherConfig, processor batchprocessor.BatchProcessor) *MicroBatcher {
	scheduler := NewBatchScheduler(config, processor)
	return &MicroBatcher{
		scheduler: scheduler,
	}
}

func (mb *MicroBatcher) Submit(j *job.Job) <-chan job.JobResult {
	mb.scheduler.Schedule(j)
	return j.ResultChan()
}

func (mb *MicroBatcher) Start() {
	log.Printf("Starting microbatcher")
	mb.scheduler.StartProcessing()
}

func (mb *MicroBatcher) Shutdown() {
	log.Printf("Shutting down microbatcher")
	mb.scheduler.Stop()
	jobqueue.DrainQueue()
}
