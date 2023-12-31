package server

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/1x-eng/mbat/pkg/jobqueue"
	"github.com/1x-eng/mbat/pkg/microbatcher"
	"github.com/gofiber/fiber/v2"
)

func NewMicrobatchingServer(batcher *microbatcher.MicroBatcher) *fiber.App {
	app := fiber.New()
	setupRoutes(app, batcher)
	return app
}

func fiberApp(app *fiber.App, port string) {
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func worker(batcher *microbatcher.MicroBatcher) {
	for job := range jobqueue.JobQueue {
		batcher.Submit(job)
	}
}

func Start(app *fiber.App, batcher *microbatcher.MicroBatcher, port string, queueSize int, processedJobsCacheTTL time.Duration, processedJobsCacheCleanup time.Duration) {
	jobqueue.InitQueue(queueSize)
	jobqueue.InitProcessedJobsCache(processedJobsCacheTTL, processedJobsCacheCleanup)

	go fiberApp(app, port)
	go worker(batcher)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	batcher.Shutdown()
	log.Println("Server and batcher shut down gracefully")
}
