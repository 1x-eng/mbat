package microbatcher

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/1x-eng/mbat/pkg/batchprocessor"
	"github.com/1x-eng/mbat/pkg/job"
	"github.com/1x-eng/mbat/pkg/jobqueue"
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

func (bs *BatchScheduler) StartProcessing() {
	log.Printf("Starting processing with batch size %d and interval %s\n", bs.config.BatchSize, bs.config.BatchInterval.String())
	bs.wg.Add(1)
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
		jobqueue.StoreProcessedJob(batch[i])
		jobqueue.Dequeue(batch[i].ID)
		log.Printf("Processed job %s. Job persisted to processed jobs cache & dequeued\n", batch[i].ID)

	}
	log.Println("Batch processing complete")
}

func (bs *BatchScheduler) processAnyRemainingJobs(batch *[]*job.Job) {
	if len(*batch) > 0 {
		log.Printf("Processing remaining %d jobs in batch\n", len(*batch))
		bs.processBatch(*batch)
		*batch = nil // Clear the batch
	} else {
		log.Println("No jobs to process in batch")
	}
}

func (bs *BatchScheduler) Stop() {
	log.Println("Sending quit signal to batching routine")
	close(bs.quit)

	bs.wg.Wait()
	log.Println("Batch scheduler wait group finished")

	close(bs.jobs)
	log.Println("Jobs channel closed")
}

func (bs *BatchScheduler) addToBatchAndProcessIfFull(job *job.Job, batch *[]*job.Job, batchTimer *time.Timer) {
	*batch = append(*batch, job)
	if len(*batch) >= bs.config.BatchSize {
		log.Printf("Batch full, processing batch of %d jobs\n", len(*batch))
		bs.processBatch(*batch)

		// Clear the batch & reset interval.
		*batch = nil
		if !batchTimer.Stop() {
			log.Println("Batch timer not stopped, draining channel")
			<-batchTimer.C // Drain the timer channel if it wasn't stopped
		}
		batchTimer.Reset(bs.config.BatchInterval)
		log.Println("A full batch was processed, resetting batch timer")
	}
}

func (bs *BatchScheduler) batchingRoutine() {
	defer bs.wg.Done()

	log.Println("Batching routine started")
	var batch []*job.Job
	var batchTimer = time.NewTimer(bs.config.BatchInterval)

	for {
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
			log.Println("Remaining jobs in the batch, if any - processed, resetting batch timer")
			batchTimer.Reset(bs.config.BatchInterval)

		case <-bs.quit:
			log.Println("Quit signal received, processing remaining jobs and stopping")
			bs.processAnyRemainingJobs(&batch)
			return
		}
	}
}
