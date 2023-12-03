package jobqueue

import (
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
}

func GetProcessedJobByID(jobID string) (*job.Job, bool) {
	j, found := processedJobsCache.Get(jobID)
	if found {
		return j.(*job.Job), true
	}
	return nil, false
}
