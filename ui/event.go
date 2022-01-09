package ui

import (
	"context"
	"github.com/go-ping/ping"
	"github.com/google/uuid"
	"github.com/pterm/pterm"
	"github.com/soheltarir/ekko/config"
	"github.com/soheltarir/ekko/logger"
	"go.uber.org/zap"
)

type Element int

const (
	StatusInfo Element = iota
	TableRow
)

type Event struct {
	ID      uuid.UUID
	Element Element
	Data    interface{}
}

func NewStatusEvent(status config.ConsumerStatus) Event {
	return Event{
		ID:      uuid.New(),
		Element: StatusInfo,
		Data:    status,
	}
}

func NewTableRowEvent(server config.Server, stats *ping.Statistics, err error) Event {
	data := map[string]interface{}{
		"server": server,
		"stats":  stats,
	}
	if err != nil {
		data["error"] = err.Error()
	} else {
		data["error"] = ""
	}
	return Event{Element: TableRow, Data: data}
}

func (u EkkoUI) Listen(ctx context.Context) {
	if !config.Config.UIEnabled {
		// Skip if UI is disabled
		return
	}

	for {
		select {
		case event := <-u.eventChan:
			if event.Element == StatusInfo {
				u.consumerStatus = event.Data.(config.ConsumerStatus)
			} else if event.Element == TableRow {
				// Build the row for the event
				eventData := event.Data.(map[string]interface{})
				server := eventData["server"].(config.Server)
				stats := eventData["stats"].(*ping.Statistics)
				err := eventData["error"].(string)
				row := statsRow(server, stats, err)
				// Update the corresponding row
				u.rows[u.addressMap[server.Address]] = row
			}
			u.render()
		case <-ctx.Done():
			logger.Log.Warn("Received termination signal, stopping listening to changes in UI")
			return
		}
	}
}

func (u EkkoUI) clearDisplay() {
	print("\033[H\033[2J")
}

func (u EkkoUI) render() {
	table, err := networkStatsTable(u.rows)
	if err != nil {
		// Skip render
		logger.Log.Warn("Failed to render network stats table", zap.Error(err))
		return
	}
	panels := pterm.Panels{
		// Header panel
		{{Data: header()}},
		// New Line
		{{Data: pterm.Sprintln()}},
		// Program information
		{{Data: pterm.DefaultBasicText.Sprint(programInfo(u.consumerStatus)...)}},
		// New line
		{{Data: pterm.Sprintln()}},
		// Network Stats table
		{{Data: table}},
	}
	u.clearDisplay()
	if err := pterm.DefaultPanel.WithPanels(panels).Render(); err != nil {
		logger.Log.Panic("Failed to render panels", zap.Error(err))
	}
}
