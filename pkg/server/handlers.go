package server

import (
	"github.com/1x-eng/mbat/pkg/job"
	"github.com/1x-eng/mbat/pkg/jobqueue"
	"github.com/1x-eng/mbat/pkg/microbatcher"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

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
