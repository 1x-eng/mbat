package microbatcher

import (
	"sync"
	"time"

	"github.com/1x-eng/mbat/pkg/batchprocessor"
	"github.com/1x-eng/mbat/pkg/job"
)

type MicroBatcherConfig struct {
	BatchSize     int
	BatchInterval time.Duration
}

type MicroBatcher struct {
	scheduler *BatchScheduler
	wg        sync.WaitGroup
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
	mb.scheduler.StartProcessing(&mb.wg)
}

func (mb *MicroBatcher) Shutdown() {
	mb.scheduler.Stop()
	mb.wg.Wait()
}
