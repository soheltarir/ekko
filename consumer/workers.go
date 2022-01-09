package consumer

import (
	"context"
	"github.com/soheltarir/ekko/config"
	"github.com/soheltarir/ekko/logger"
	"go.uber.org/zap"
	"sync"
)

// Start acts as the proxy between the ingestChan and jobsChan,
// with a select to support graceful shutdown.
func (c *Consumer) Start(ctx context.Context) {
	c.Status = config.Running
	c.notifyStatusRenderer(c.Status)
	for {
		select {
		case <-ctx.Done():
			c.HandleShutdown()
			return
		case job := <-c.channels.ingestion:
			logger.Log.Debug("Received job", zap.Any("event", job))
			c.channels.job <- job
		}
	}
}

// SpawnWorker starts a thread which listens for ping job events, and executes them
func (c *Consumer) SpawnWorker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	log := logger.Log.With(zap.Int("worker_id", id))
	log.Debug("Ping worker started")

	for job := range c.channels.job {
		// Run the ping for received event
		c.ping(job.Destination, log)
	}
	log.Warn("Shutdown signal received, stopping worker...")
}
