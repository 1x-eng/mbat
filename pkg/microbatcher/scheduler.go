package microbatcher

import (
	"fmt"
	"log"
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
		log.Println("Job successfully scheduled")
	default:
		errorMsg := "scheduler is either full or not accepting new jobs"
		log.Println(errorMsg)
		j.SetResult(job.JobResult{Error: fmt.Errorf(errorMsg)})
	}
}

func (bs *BatchScheduler) StartProcessing(wg *sync.WaitGroup) {
	wg.Add(1)
	go bs.batchingRoutine()
}

func (bs *BatchScheduler) processBatch(batch []*job.Job) {
	if len(batch) == 0 {
		log.Println("No jobs to process in batch")
		return
	}

	log.Printf("Processing batch of %d jobs\n", len(batch))
	results := bs.batchProcessor.ProcessBatch(batch)
	for i, result := range results {
		if result.Error != nil {
			log.Printf("Error processing job: %v\n", result.Error)
		}
		batch[i].SetResult(result)
	}
	log.Println("Batch processing complete")
}

func (bs *BatchScheduler) processAnyRemainingJobs(batch *[]*job.Job) {
	if len(*batch) > 0 {
		bs.processBatch(*batch)
		*batch = nil // Clear the batch
	}
}

func (bs *BatchScheduler) Stop() {
	close(bs.quit)
	bs.wg.Wait()
	close(bs.jobs)
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

	log.Println("Batching routine started")
	var batch []*job.Job
	var batchTimer *time.Timer

	for {
		batchTimer = bs.newBatchTimer(batchTimer, bs.config.BatchInterval)

		select {
		case job, ok := <-bs.jobs:
			if !ok {
				log.Println("Jobs channel closed, processing remaining jobs")
				bs.processAnyRemainingJobs(&batch)
				return
			}
			log.Println("Job received, adding to batch")
			bs.addToBatchAndProcessIfFull(job, &batch, batchTimer)

		case <-batchTimer.C:
			log.Println("Batch timer expired, processing any jobs in batch")
			bs.processAnyRemainingJobs(&batch)

		case <-bs.quit:
			log.Println("Quit signal received, processing remaining jobs and stopping")
			bs.processAnyRemainingJobs(&batch)
			return
		}
	}
}
