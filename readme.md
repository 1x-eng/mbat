# µ-batching 

## Overview
A simple microbatching system in Go, designed to process jobs efficiently in batches. It includes a queue for job management and a batch processing mechanism. The system also provides a REST API to submit jobs, check job status, and retrieve information about queued and processed jobs.

## Features
- Job Submission: Submit jobs for processing via a REST API.
- Job Queueing: Jobs are queued before being batch processed.
- Batch Processing: Jobs are processed in configurable batch sizes and intervals.
- Job Status Tracking: Track the status of submitted jobs, whether they're queued or processed.
- Job ID Retrieval: Retrieve the IDs of all jobs currently in the queue.
- Processed Job Retrieval: Retrieve all processed jobs.

## Getting Started

### Prerequisites
Go (version 1.15 or later)

### Installation & hey, just run it
Clone the repository to your local machine:

```bash
git clone https://github.com/1x-eng/mbat.git
cd mbat
touch ./.env # edit the config as you see fit

# build and run the executable
go build -o ./bin/mbat ./cmd/
./bin/mbat
```

## API Endpoints
- Submit Job: `POST /job/submit`
- Job Status: `GET /job/status/:jobID`
- Queued Jobs: `GET /jobs/queue`
- Processed Jobs: `GET /jobs/processed`

## Swagger / OpenAPI
- µ-batching uses `swaggo` cli to generate swagger specs. Install swaggo cli using - `go install github.com/swaggo/swag/cmd/swag@latest`
- To generate or refresh swagger specs, `swag init -g ./pkg/server/handlers.go`
- µ-batching also uses fiber-swagger to render swaghger UI.
- Once the project is built and run, visit `<api>/swagger` to view swagger UI.

## TODO
- Integrate github actions / promote CI/CD and abstract swaggo / fiber-swagger
- A whole lot of tests :P
