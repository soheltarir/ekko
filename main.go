package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	Log.Info("Ekko service started")
	PrintIntro()

	// Set up cancellation context and wait group
	ctx, cancelFunc := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	// create the consumer
	consumer := NewConsumer()

	// Send the servers to ping as events to worker/s
	producer := Producer{callbackFunc: consumer.callbackFunc}
	go producer.start(ctx)

	// Render table and listen for updates
	table := NewStatsTable(Config.Servers, consumer.UIEventChan)
	go table.listenForChange(ctx)

	// Start consumer with cancellation context passed
	go consumer.startConsumer(ctx)

	// Start workers and Add [workerPoolSize] to WaitGroup
	wg.Add(Config.WorkerPoolSize)
	for i := 0; i < Config.WorkerPoolSize; i++ {
		go consumer.worker(wg, i)
	}

	// Handle sigterm and await termChan signal
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-termChan // Blocks here until interrupted

	// Handle shutdown
	Log.Warn("Shutdown signal received")
	cancelFunc() // Signal cancellation to context.Context
	wg.Wait()    // Block here until are workers are done
	consumer.closeUIChannel()
	Log.Debug("All workers stopped, shutting down")
}
