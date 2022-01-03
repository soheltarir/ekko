package main

import (
	"context"
	"github.com/soheltarir/ekko/config"
	"github.com/soheltarir/ekko/consumer"
	"github.com/soheltarir/ekko/logger"
	"github.com/soheltarir/ekko/ui"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	logger.Log.Info("Ekko service started")

	// Set up cancellation context and wait group
	ctx, cancelFunc := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	// Create an event for UI changes
	uiChan := make(chan ui.Event)

	// Render UI and listen for updates
	ekkoUI := ui.New(config.Config.Servers, uiChan)
	go ekkoUI.Listen(ctx)

	// create the consumer
	pingConsumer := consumer.New(uiChan)

	// Start consumer with cancellation context passed
	go pingConsumer.Start(ctx)

	// Start workers and Add [workerPoolSize] to WaitGroup
	wg.Add(config.Config.WorkerPoolSize)
	for i := 0; i < config.Config.WorkerPoolSize; i++ {
		go pingConsumer.SpawnWorker(i, wg)
	}

	// Send the servers to ping as events to worker/s
	producer := Producer{callbackFunc: pingConsumer.CallbackFunc}
	go producer.Start(ctx)

	// Handle sigterm and await termChan signal
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-termChan // Blocks here until interrupted

	// Handle shutdown
	logger.Log.Warn("Shutdown signal received")
	cancelFunc() // Signal cancellation to context.Context
	wg.Wait()    // Block here until are workers are done
	logger.Log.Debug("All workers stopped, shutting down")
}
