package consumer

import (
	"github.com/google/uuid"
	"github.com/soheltarir/ekko/config"
	"github.com/soheltarir/ekko/ui"
	"sync"
	"time"
)

// Event defines the object a ping job received
type Event struct {
	ID          uuid.UUID
	Destination config.Server
	SentOn      time.Time
}

// NewEvent creates a new ping job
func NewEvent(dest config.Server) Event {
	return Event{
		ID:          uuid.New(),
		Destination: dest,
		SentOn:      time.Now(),
	}
}

// channels object encapsulates all go channels that would be used to
// communicate between consumers & producers
type channels struct {
	// ingestion channel is used as a proxy between producer and ping consumers
	ingestion chan Event
	// job channel is used to queue ping job to a Destination
	job chan Event
	// ui channel is used to queue UI updates and further consumed by UI renderer
	ui chan ui.Event
}

// Consumer exposes methods and parameters to control the behaviour of the ping workers
type Consumer struct {
	channels channels
	// activeJobs contains the list of actively running ping jobs
	activeJobs sync.Map
	lock       sync.Mutex
	// Status defines the current running state of the consumer
	Status config.ConsumerStatus
}

// New returns a new Consumer object
func New(uiEventChan chan ui.Event) *Consumer {
	channels := channels{
		ingestion: make(chan Event),
		job:       make(chan Event),
		ui:        uiEventChan,
	}
	return &Consumer{
		channels: channels,
		Status:   config.NotStarted,
	}
}
