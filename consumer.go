package main

import (
	"context"
	"github.com/go-ping/ping"
	"go.uber.org/zap"
	"math/rand"
	"sync"
	"time"
)

// Consumer object contains channels to pass events
type Consumer struct {
	ingestChan chan Server
	jobsChan   chan Server
}

// callbackFunc is invoked each time the producer sends an event.
func (c Consumer) callbackFunc(event Server) {
	c.ingestChan <- event
}

// startConsumer acts as the proxy between the ingestChan and jobsChan,
// with a select to support graceful shutdown.
func (c Consumer) startConsumer(ctx context.Context) {
	for {
		select {
		case job := <-c.ingestChan:
			c.jobsChan <- job
		case <-ctx.Done():
			Log.Warn("Consumer received cancellation signal, closing jobs channel")
			close(c.jobsChan)
			Log.Info("Jobs channel successfully closed")
			return
		}
	}
}

func (c Consumer) workerFunc(wg *sync.WaitGroup, index int) {
	defer wg.Done()
	logger := Log.With(zap.Int("worker_id", index))
	logger.Info("Ping worker starting")
	// Create a worker specific pinger object
	pinger, err := ping.NewPinger("localhost")
	if err != nil {
		logger.Panic("Failed to initialise pinger for the worker", zap.Int("worker_id", index), zap.Error(err))
	}

	for job := range c.jobsChan {
		pingServer(job, pinger, logger)
		logger.Debug("Ping complete, putting worker to sleep")
		time.Sleep(time.Second * 20)
	}
	logger.Warn("Interrupt signal received, stopping worker")
	pinger.Stop()
}

func pingServer(server Server, pinger *ping.Pinger, logger *zap.Logger) {
	log := logger.With(
		zap.String("server_name", server.Name),
		zap.String("game", server.Game),
		zap.String("server_ip", server.IPAddress),
	)

	if err := pinger.SetAddr(server.IPAddress); err != nil {
		log.Error("Failed to initialise ping", zap.Error(err))
		return
	}

	// Randomize the count of packets to be sent
	rand.Seed(time.Now().UnixNano())
	pinger.Count = rand.Intn(Config.MaxPacketNum-Config.MinPacketNum) + Config.MinPacketNum
	// Timeout of 30
	pinger.Timeout = time.Second * time.Duration(Config.PingTimeout)

	pinger.OnSetup = func() {
		log.Info("Pinging the server")
	}

	pinger.OnFinish = func(s *ping.Statistics) {
		log.Info("Ping complete",
			zap.Int("num_packets", s.PacketsSent),
			zap.Float64("packet_loss", s.PacketLoss),
			zap.Duration("avg_rtt", s.AvgRtt),
			zap.Duration("min_rtt", s.MinRtt),
			zap.Duration("max_rtt", s.MaxRtt),
		)
	}

	if err := pinger.Run(); err != nil {
		log.Error("Failed to run ping", zap.Error(err))
		return
	}
}
