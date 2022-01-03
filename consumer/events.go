package consumer

import (
	"github.com/go-ping/ping"
	"github.com/soheltarir/ekko/config"
	"github.com/soheltarir/ekko/logger"
	"github.com/soheltarir/ekko/ui"
	"go.uber.org/zap"
)

// CallbackFunc is invoked each time the producer sends an event
func (c *Consumer) CallbackFunc(event Event) {
	logger.Log.Debug("Attempting to send event to ingestion channel", zap.Any("event", event))
	c.channels.ingestion <- event
	logger.Log.Debug("Sent event to ingestion channel", zap.Any("event", event))
}

// notifyTableRenderer sends an event to the table renderer with new data
func (c *Consumer) notifyTableRenderer(dest config.Server, stats *ping.Statistics, err error) {
	// Skip sending an event if UI is disabled or the consumer is not running
	if !config.Config.UIEnabled || c.Status != config.Running {
		return
	}

	// acquire lock
	c.lock.Lock()
	defer c.lock.Unlock()

	// create a new event
	event := ui.NewTableRowEvent(dest, stats, err)

	logger.Log.Debug("Attempting to send UI event", zap.Any("event", event))
	c.channels.ui <- event
	logger.Log.Debug("Successfully sent UI event", zap.Any("key", event))
}

func (c *Consumer) notifyStatusRenderer(status config.ConsumerStatus) {
	// Skip sending an event if UI is disabled or the consumer is not running
	if !config.Config.UIEnabled || c.Status != config.Running {
		return
	}

	// acquire lock
	c.lock.Lock()
	defer c.lock.Unlock()

	// create a new event
	event := ui.NewStatusEvent(status)

	logger.Log.Debug("Attempting to send UI event", zap.Any("event", event))
	c.channels.ui <- event
	logger.Log.Debug("Successfully sent UI event", zap.Any("event", event))

}
