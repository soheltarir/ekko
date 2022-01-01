package main

import (
	"context"
	"go.uber.org/zap"
	"time"
)

// Producer invokes the consumer callback function and sends a destination as an event.
// the producer repeats sending of events after sleeping for Config.PingInterval seconds,
// i.e., the destinations are pinged every 30 Config.PingInterval seconds.
type Producer struct {
	callbackFunc func(event Server)
}

// start runs the producer to trigger events indefinitely
func (p Producer) start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			Log.Debug("Producer received cancellation signal, exiting...")
			return
		default:
			for _, server := range Config.Servers {
				Log.Debug("Sending event", zap.String("name", server.Name))
				p.callbackFunc(server)
			}
			Log.Debug("Producer sleeping", zap.Int64("sleep_duration", Config.PingInterval))
			time.Sleep(time.Duration(Config.PingInterval) * time.Second)
		}
	}
}
