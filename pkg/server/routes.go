package server

import (
	_ "github.com/1x-eng/mbat/docs"
	"github.com/1x-eng/mbat/pkg/microbatcher"
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func setupRoutes(app *fiber.App, batcher *microbatcher.MicroBatcher) {
	app.Get("/swagger/*", fiberSwagger.WrapHandler)
	app.Post("/job/submit", jobSubmitHandler(batcher))
	app.Get("/job/status/:jobID", jobStatusHandler())
	app.Get("/jobs/queued", queuedJobsHandler())
	app.Get("/jobs/processed", processedJobsHandler())
}
