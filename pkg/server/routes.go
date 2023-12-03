package server

import (
	"github.com/1x-eng/mbat/pkg/microbatcher"
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App, batcher *microbatcher.MicroBatcher) {
	app.Post("/job/submit", jobSubmitHandler(batcher))
	app.Get("/job/status/:jobID", jobStatusHandler())
	app.Get("/jobs/queued", queuedJobsHandler())
	app.Get("/jobs/processed", processedJobsHandler())
}
