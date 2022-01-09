package consumer

import (
	"github.com/go-ping/ping"
	"github.com/soheltarir/ekko/config"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

func (c *Consumer) ping(destination config.Server, log *zap.Logger) {
	// Initialise a base logger
	log = log.With(
		zap.String("server_name", destination.Name),
		zap.String("server_ip", destination.Address),
		zap.Any("labels", destination.Labels),
	)
	pingThread, err := ping.NewPinger(destination.Address)
	if err != nil {
		log.Error("Failed to initialise ping", zap.Error(err))
		c.notifyTableRenderer(destination, nil, err)
		return
	}

	// Randomize the count of packets to be sent
	rand.Seed(time.Now().UnixNano())
	pingThread.Count = rand.Intn(config.Config.MaxPacketNum-config.Config.MinPacketNum) + config.Config.MinPacketNum

	// Set the timeout for a packet to consider it as failed
	pingThread.Timeout = time.Second * time.Duration(config.Config.PingTimeout)

	// Override the default logger
	pingThread.SetLogger(log.Sugar())

	// Run as privileged user to promote connections to ICMP
	pingThread.SetPrivileged(true)

	// Initialise a callback function to run when the ping starts
	pingThread.OnSetup = func() {
		log.Info("Ping started")
		// Add the ping job to active list
		c.activeJobs.Store(pingThread.ID(), pingThread)
	}

	// Initialise a callback function to run when the ping finishes
	pingThread.OnFinish = func(s *ping.Statistics) {
		log.Info("Ping complete",
			zap.Int("num_packets", s.PacketsSent),
			zap.Float64("packet_loss", s.PacketLoss),
			zap.Duration("avg_rtt", s.AvgRtt),
			zap.Duration("min_rtt", s.MinRtt),
			zap.Duration("max_rtt", s.MaxRtt),
		)
		c.notifyTableRenderer(destination, s, nil)
		// Delete the ping job from active list
		c.activeJobs.Delete(pingThread.ID())
	}

	if err := pingThread.Run(); err != nil {
		log.Error("Failed to run ping", zap.Error(err))
		c.notifyTableRenderer(destination, nil, err)
		return
	}
}
