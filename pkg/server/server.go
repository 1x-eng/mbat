package server

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/1x-eng/mbat/pkg/jobqueue"
	"github.com/1x-eng/mbat/pkg/microbatcher"
	"github.com/gofiber/fiber/v2"
)

func NewMicrobatchingServer(batcher *microbatcher.MicroBatcher) *fiber.App {
	app := fiber.New()
	setupRoutes(app, batcher)
	return app
}

func Start(app *fiber.App, batcher *microbatcher.MicroBatcher, port string, queueSize int) {
	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
		jobqueue.InitQueue(queueSize)
		log.Printf("Microbatching server started on port %s", port)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	batcher.Shutdown()
	log.Println("Server and batcher shut down gracefully")
}
