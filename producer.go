package main

import "time"

// Producer simulates an external library that invokes the
// registered callback when it has new data for us once per 100ms.
type Producer struct {
	callbackFunc func(event Server)
}

func (p Producer) start() {
	currItr, size := 0, len(Config.Servers)
	for {
		p.callbackFunc(Config.Servers[currItr])
		currItr++
		if currItr >= size {
			currItr = 0
			time.Sleep(time.Minute)
		}
	}
}
