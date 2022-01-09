package main

import (
	"context"
	"github.com/soheltarir/ekko/config"
	"github.com/soheltarir/ekko/consumer"
	"github.com/soheltarir/ekko/logger"
	"go.uber.org/zap"
	"time"
)

// Producer invokes the consumer callback function and sends a destination as an event.
// the producer repeats sending of events after sleeping for Config.PingInterval seconds,
// i.e., the destinations are pinged every 30 Config.PingInterval seconds.
type Producer struct {
	callbackFunc func(event consumer.Event)
}

// Start runs the producer to trigger events indefinitely
func (p Producer) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logger.Log.Debug("Producer received cancellation signal, exiting...")
			return
		default:
			for _, server := range config.Config.Servers {
				pingEvent := consumer.NewEvent(server)
				logger.Log.Debug("Sending event", zap.Any("event", pingEvent))
				p.callbackFunc(pingEvent)
			}
			logger.Log.Debug("Producer sleeping", zap.Int64("sleep_duration", config.Config.PingInterval))
			time.Sleep(time.Duration(config.Config.PingInterval) * time.Second)
		}
	}
}
