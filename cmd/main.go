package main

import (
	"log"

	"github.com/1x-eng/mbat/pkg/config"
	"github.com/1x-eng/mbat/pkg/microbatcher"
	"github.com/1x-eng/mbat/pkg/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	mbatcher := microbatcher.NewMicroBatcher(cfg.MicroBatcherConfig, cfg.Processor)
	mbatcher.Start()

	app := server.NewMicrobatchingServer(mbatcher)
	server.Start(app, mbatcher, cfg.Port, cfg.QueueSize)
}
