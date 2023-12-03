package config

import (
	"os"
	"strconv"
	"time"

	"github.com/1x-eng/mbat/pkg/batchprocessor"
	"github.com/1x-eng/mbat/pkg/microbatcher"
	"github.com/joho/godotenv"
)

type Config struct {
	MicroBatcherConfig        microbatcher.MicroBatcherConfig
	Processor                 batchprocessor.BatchProcessor
	Port                      string
	QueueSize                 int
	ProcessedJobsCacheTTL     time.Duration
	ProcessedJobsCacheCleanup time.Duration
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	batchSize, err := strconv.Atoi(os.Getenv("BATCH_SIZE"))
	if err != nil {
		return nil, err
	}

	batchIntervalSeconds, err := strconv.Atoi(os.Getenv("BATCH_INTERVAL_SECONDS"))
	if err != nil {
		return nil, err
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	queueSize, err := strconv.Atoi(os.Getenv("QUEUE_SIZE"))
	if err != nil {
		queueSize = 100
	}

	processedJobsCacheTTLSeconds, err := strconv.Atoi(os.Getenv("PROCESSED_JOBS_CACHE_TTL"))
	if err != nil {
		processedJobsCacheTTLSeconds = 3600 // 1 hr.
	}

	processedJobCleanupSeconds, err := strconv.Atoi(os.Getenv("PROCESSED_JOBS_CACHE_CLEANUP_MINUTES"))
	if err != nil {
		processedJobCleanupSeconds = 600 // 10 mins.
	}

	return &Config{
		MicroBatcherConfig: microbatcher.MicroBatcherConfig{
			BatchSize:     batchSize,
			BatchInterval: time.Duration(batchIntervalSeconds) * time.Second,
		},
		Processor:                 batchprocessor.NewMockBatchProcessor(),
		Port:                      port,
		QueueSize:                 queueSize,
		ProcessedJobsCacheTTL:     time.Duration(processedJobsCacheTTLSeconds) * time.Second,
		ProcessedJobsCacheCleanup: time.Duration(processedJobCleanupSeconds) * time.Second,
	}, nil
}
