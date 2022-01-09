package consumer

import (
	"github.com/go-ping/ping"
	"github.com/soheltarir/ekko/config"
	"github.com/soheltarir/ekko/logger"
	"go.uber.org/zap"
)

// closeUiChannel closes the ui channel to stop receiving any further events
func (c *Consumer) closeUiChannel() {
	c.lock.Lock()
	defer c.lock.Unlock()

	logger.Log.Warn("Closing UI events channel")
	close(c.channels.ui)
	logger.Log.Debug("UI events channel successfully closed")
}

// stopActiveJobs attempts to stop all actively running ping jobs
func (c *Consumer) stopActiveJobs() {
	c.activeJobs.Range(func(key, value interface{}) bool {
		job := value.(*ping.Pinger)
		logger.Log.Warn("Stopping ping job", zap.Int("job_id", key.(int)),
			zap.String("address", job.Addr()))
		job.Stop()
		logger.Log.Debug("Successfully stopped job", zap.Int("job_id", key.(int)),
			zap.String("address", job.Addr()))
		return true
	})
}

// HandleShutdown method triggers all stop instructions when a shutdown signal is received
func (c *Consumer) HandleShutdown() {
	c.Status = config.Stopped
	logger.Log.Warn("Consumer received cancellation signal, closing jobs channel")
	close(c.channels.job)
	logger.Log.Debug("Jobs channel successfully closed")
	c.stopActiveJobs()
	c.notifyStatusRenderer(c.Status)
	c.closeUiChannel()
}
