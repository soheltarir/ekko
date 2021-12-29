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
	// activePingJobs contains the list of actively running ping jobs
	activePingJobs sync.Map
	lock           sync.Mutex
	// uiEventChan is the channel via which UI change events will be sent
	uiEventChan chan UIEvent
}

// callbackFunc is invoked each time the producer sends an event.
func (c *Consumer) callbackFunc(event Server) {
	c.ingestChan <- event
}

// notifyTableRenderer sends an event to the table renderer with new data
func (c *Consumer) notifyTableRenderer(dest Server, stats *ping.Statistics, err string) {
	c.lock.Lock()
	c.uiEventChan <- UIEvent{dest: dest, stats: stats, err: err}
	c.lock.Unlock()
}

// stopActivePingJobs attempts to stop all actively running ping jobs
func (c *Consumer) stopActivePingJobs() {
	c.activePingJobs.Range(func(key, value interface{}) bool {
		pinger := value.(*ping.Pinger)
		Log.Info("Stopping pinger", zap.Int("pinger_id", key.(int)),
			zap.String("address", pinger.Addr()))
		pinger.Stop()
		Log.Info("Successfully stopped pinger", zap.Int("pinger_id", key.(int)),
			zap.String("address", pinger.Addr()))
		return true
	})
}

// startConsumer acts as the proxy between the ingestChan and jobsChan,
// with a select to support graceful shutdown.
func (c *Consumer) startConsumer(ctx context.Context) {
	for {
		select {
		case job := <-c.ingestChan:
			c.jobsChan <- job
		case <-ctx.Done():
			Log.Warn("Consumer received cancellation signal, closing jobs channel")
			close(c.jobsChan)
			Log.Info("Jobs channel successfully closed")
			c.stopActivePingJobs()
			return
		}
	}
}

// worker starts a thread which listens for ping job events, and executes them
func (c *Consumer) worker(wg *sync.WaitGroup, index int) {
	defer wg.Done()
	logger := Log.With(zap.Int("worker_id", index))
	logger.Info("Ping worker starting")

	for job := range c.jobsChan {
		c.pingServer(job, logger)
	}
	logger.Warn("Interrupt signal received, stopping worker")
}

// pingServer sends ICMP/UDP packets to the destination and records network statistics
func (c *Consumer) pingServer(server Server, logger *zap.Logger) {
	log := logger.With(
		zap.String("server_name", server.Name),
		zap.String("server_ip", server.Address),
		zap.Any("labels", server.Labels),
	)
	pinger, err := ping.NewPinger(server.Address)
	if err != nil {
		log.Error("Failed to initialise ping", zap.Error(err))
		c.notifyTableRenderer(server, nil, err.Error())
		return
	}

	// Randomize the count of packets to be sent
	rand.Seed(time.Now().UnixNano())
	pinger.Count = rand.Intn(Config.MaxPacketNum-Config.MinPacketNum) + Config.MinPacketNum
	// Set the timeout for a packet to consider it as failed
	pinger.Timeout = time.Second * time.Duration(Config.PingTimeout)
	// Override the default logger
	pinger.SetLogger(log.Sugar())

	pinger.OnSetup = func() {
		log.Info("Ping started")
		// Add the ping job to active list
		c.activePingJobs.Store(pinger.ID(), pinger)
	}

	pinger.OnFinish = func(s *ping.Statistics) {
		log.Info("Ping complete",
			zap.Int("num_packets", s.PacketsSent),
			zap.Float64("packet_loss", s.PacketLoss),
			zap.Duration("avg_rtt", s.AvgRtt),
			zap.Duration("min_rtt", s.MinRtt),
			zap.Duration("max_rtt", s.MaxRtt),
		)
		// Delete the ping job from active list
		c.activePingJobs.Delete(pinger.ID())
		c.notifyTableRenderer(server, s, "")
	}

	if err := pinger.Run(); err != nil {
		log.Error("Failed to run ping", zap.Error(err))
		c.notifyTableRenderer(server, nil, err.Error())
		return
	}
}

// NewConsumer constructs a consumer object which runs workers to ping the specified destinations
func NewConsumer(uiEventChan chan UIEvent) *Consumer {
	return &Consumer{
		ingestChan:  make(chan Server),
		jobsChan:    make(chan Server),
		uiEventChan: uiEventChan,
	}
}
