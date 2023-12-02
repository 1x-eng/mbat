package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/1x-eng/mbat/pkg/batchprocessor"
	"github.com/1x-eng/mbat/pkg/job"
	"github.com/1x-eng/mbat/pkg/microbatcher"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	batchSize, err := strconv.Atoi(os.Getenv("BATCH_SIZE"))
	if err != nil {
		log.Fatal("Invalid BATCH_SIZE value")
	}
	batchIntervalSeconds, err := strconv.Atoi(os.Getenv("BATCH_INTERVAL_SECONDS"))
	if err != nil {
		log.Fatal("Invalid BATCH_INTERVAL_SECONDS value")
	}

	config := microbatcher.MicroBatcherConfig{
		BatchSize:     batchSize,
		BatchInterval: time.Duration(batchIntervalSeconds) * time.Second,
	}

	processor := batchprocessor.NewMockBatchProcessor() // DI batch processor here

	batcher := microbatcher.NewMicroBatcher(config, processor)
	batcher.Start()

	// todo: expose as a cli arg or maybe make a rest api?
	j := job.NewJob("example1: data")
	resultChan := batcher.Submit(j)

	result := <-resultChan
	if result.Error != nil {
		log.Printf("Error processing job: %v\n", result.Error)
	} else {
		log.Printf("Job processed successfully: %v\n", result.Data)
	}

	batcher.Shutdown()
}
