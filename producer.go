package main

import "time"

// Producer invokes the consumer callback function and sends a destination as an event.
// the producer repeats sending of events after sleeping for Config.PingInterval seconds,
// i.e., the destinations are pinged every 30 Config.PingInterval seconds.
type Producer struct {
	callbackFunc func(event Server)
}

// start runs the producer to trigger events indefinitely
func (p Producer) start() {
	for {
		for _, server := range Config.Servers {
			p.callbackFunc(server)
		}
		time.Sleep(time.Duration(Config.PingInterval) * time.Second)
	}
}
