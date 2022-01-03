package ui

import (
	"github.com/go-ping/ping"
	"github.com/soheltarir/ekko/config"
)

// EkkoUI exposes all methods and objects for displaying network statistics
type EkkoUI struct {
	// rows is the list of destination stats (including the header)
	rows [][]string
	// addressMap stores the link between the destination address, and it's
	// corresponding row index
	addressMap map[string]int
	// eventChan is used for listening to UI change events
	eventChan      chan Event
	consumerStatus config.ConsumerStatus
}

func New(destinations []config.Server, uiChan chan Event) *EkkoUI {
	if !config.Config.UIEnabled {
		// Skip if UI is disabled
		return &EkkoUI{}
	}
	ui := EkkoUI{
		rows:           make([][]string, 0),
		addressMap:     make(map[string]int),
		eventChan:      uiChan,
		consumerStatus: config.NotStarted,
	}
	// Add the header row
	ui.rows = append(ui.rows, StatsTableHeader)
	// Populate the ui with the destinations containing empty Data
	for idx, dest := range destinations {
		// Create an empty stats row for initialisation
		stats := &ping.Statistics{Addr: dest.Address}
		ui.rows = append(ui.rows, statsRow(dest, stats, ""))
		ui.addressMap[dest.Address] = idx + 1
	}
	ui.render()
	return &ui
}
