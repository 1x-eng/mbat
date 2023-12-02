package microbatcher

import (
	"fmt"
	"sync"
	"time"

	"github.com/1x-eng/mbat/pkg/batchprocessor"
	"github.com/1x-eng/mbat/pkg/job"
)

type BatchScheduler struct {
	config         MicroBatcherConfig
	jobs           chan *job.Job
	batchProcessor batchprocessor.BatchProcessor
	quit           chan struct{}
	wg             sync.WaitGroup
}

func NewBatchScheduler(config MicroBatcherConfig, processor batchprocessor.BatchProcessor) *BatchScheduler {
	return &BatchScheduler{
		config:         config,
		jobs:           make(chan *job.Job, config.BatchSize),
		batchProcessor: processor,
		quit:           make(chan struct{}),
	}
}

func (bs *BatchScheduler) Schedule(j *job.Job) {
	select {
	case bs.jobs <- j:
		// Job successfully scheduled
	default:
		j.SetResult(job.JobResult{Error: fmt.Errorf("scheduler is either full or not accepting new jobs")})
	}
}

func (bs *BatchScheduler) StartProcessing(wg *sync.WaitGroup) {
	wg.Add(1)
	go bs.batchingRoutine()
}

func (bs *BatchScheduler) processBatch(batch []*job.Job) {
	if len(batch) == 0 {
		return
	}

	results := bs.batchProcessor.ProcessBatch(batch)
	for i, result := range results {
		batch[i].SetResult(result)
	}
}

func (bs *BatchScheduler) processAnyRemainingJobs(batch *[]*job.Job) {
	if len(*batch) > 0 {
		bs.processBatch(*batch)
		*batch = nil // Clear the batch
	}
}

func (bs *BatchScheduler) Stop() {
	close(bs.quit)
	close(bs.jobs)
	bs.wg.Wait()
}

func (bs *BatchScheduler) addToBatchAndProcessIfFull(job *job.Job, batch *[]*job.Job, batchTimer *time.Timer) {
	*batch = append(*batch, job)
	if len(*batch) >= bs.config.BatchSize {
		bs.processBatch(*batch)

		// Clear the batch & reset interval.
		*batch = nil
	}
}

func (bs *BatchScheduler) newBatchTimer(previousTimer *time.Timer, duration time.Duration) *time.Timer {
	if previousTimer != nil {
		if !previousTimer.Stop() {
			<-previousTimer.C // Drain the channel if the timer had already fired
		}
	}
	return time.NewTimer(duration)
}

func (bs *BatchScheduler) batchingRoutine() {
	defer bs.wg.Done()

	var batch []*job.Job
	var batchTimer *time.Timer

	for {
		batchTimer = bs.newBatchTimer(batchTimer, bs.config.BatchInterval)

		select {
		case job, ok := <-bs.jobs:
			if !ok {
				bs.processAnyRemainingJobs(&batch)
				return
			}
			bs.addToBatchAndProcessIfFull(job, &batch, batchTimer)

		case <-batchTimer.C:
			bs.processAnyRemainingJobs(&batch)

		case <-bs.quit:
			bs.processAnyRemainingJobs(&batch)
			return
		}
	}
}
