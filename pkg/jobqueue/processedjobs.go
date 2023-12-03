package jobqueue

import (
	"log"
	"time"

	"github.com/1x-eng/mbat/pkg/job"
	"github.com/patrickmn/go-cache"
)

var processedJobsCache *cache.Cache

func InitProcessedJobsCache(defaultTTL time.Duration, cleanup time.Duration) {
	processedJobsCache = cache.New(defaultTTL, cleanup)
}

func StoreProcessedJob(j *job.Job) {
	processedJobsCache.Set(j.ID.String(), j, cache.DefaultExpiration)
	log.Printf("Stored job %s in processed jobs cache\n", j.ID)
}

func GetProcessedJobByID(jobID string) (*job.Job, bool) {
	j, found := processedJobsCache.Get(jobID)
	if found {
		return j.(*job.Job), true
	}
	return nil, false
}

func GetAllProcessedJobs() []*job.Job {
	var processedJobs []*job.Job
	for _, item := range processedJobsCache.Items() {
		processedJobs = append(processedJobs, item.Object.(*job.Job))
	}
	return processedJobs
}
