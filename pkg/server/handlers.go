package server

import (
	"github.com/1x-eng/mbat/pkg/job"
	"github.com/1x-eng/mbat/pkg/jobqueue"
	"github.com/1x-eng/mbat/pkg/microbatcher"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// jobSubmitHandler submits a new job, that gets enqueued & added to the batch by worker.
// @Summary Submit a new job that gets enqueued & added to the batch by worker
// @Description Expects a JSON object with a "data" field as payload, enqueues job to the queue which gets picked up by worker running in a separate goroutine; the worker adds the job to the batch. This API returns the jobID of the submitted job, which can be used to query the status of the job.
// @Accept  json
// @Produce  json
// @Param   data body string true "Job Data"
// @Success 200 {object} map[string]string
// @Router /job/submit [post]
func jobSubmitHandler(batcher *microbatcher.MicroBatcher) fiber.Handler {
	return func(c *fiber.Ctx) error {
		jData := new(struct {
			Data string `json:"data"`
		})

		if err := c.BodyParser(jData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		j := job.NewJob(jData.Data)
		jobqueue.Enqueue(j)

		return c.JSON(fiber.Map{"jobID": j.ID})
	}
}

// jobStatusHandler checks the status of a job by its ID.
// @Summary Check job status
// @Description Checks whether the job is queued or processed and returns the status along with the job result if processed.
// @Accept  json
// @Produce  json
// @Param   jobID path string true "Job ID"
// @Success 200 {object} map[string]interface{} "Job Status"
// @Router /job/status/{jobID} [get]
func jobStatusHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		jobID, err := uuid.Parse(c.Params("jobID"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid job ID"})
		}

		queuedJob, exists := jobqueue.FindJobByID(jobID.String())
		if !exists {
			processedJob, found := jobqueue.GetProcessedJobByID(jobID.String())
			if !found {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Job not found"})
			}

			return c.JSON(fiber.Map{
				"status": "processed",
				"result": processedJob.GetResult(),
			})
		}

		return c.JSON(fiber.Map{
			"status": "queued",
			"result": queuedJob.GetResult(),
		})
	}
}

// queuedJobsHandler lists all queued jobs.
// @Summary List queued jobs
// @Description Returns a list of all job IDs that are currently in the queue.
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{} "List of Queued Jobs"
// @Router /jobs/queued [get]
func queuedJobsHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		jobs := jobqueue.GetQueuedJobIds()
		return c.JSON(fiber.Map{
			"size":   len(jobs),
			"status": "queued",
			"jobs":   jobs,
		})
	}
}

// processedJobsHandler lists all processed jobs.
// @Summary List processed jobs
// @Description Returns a list of all jobs that have been processed along with their results.
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{} "List of Processed Jobs"
// @Router /jobs/processed [get]
func processedJobsHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		jobs := jobqueue.GetAllProcessedJobs()
		return c.JSON(fiber.Map{
			"size":   len(jobs),
			"status": "processed",
			"jobs":   jobs,
		})
	}
}
