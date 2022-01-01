package main

import (
	"context"
	"fmt"
	"github.com/go-ping/ping"
	"github.com/pterm/pterm"
	"go.uber.org/zap"
	"log"
	"time"
)

// StatsTableHeader contains the header row for the Network statistics table UI
var StatsTableHeader = []string{
	"Destination Name",
	"Address",
	"Packets Sent",
	"Packet Loss",
	"Avg. Response",
	"Min. Response",
	"Max. Response",
	"Time Recorded",
	"Details",
}

// UIEvent signifies an event object used for communicating between the renderer and ping jobs
type UIEvent struct {
	dest  Server
	stats *ping.Statistics
	err   string
}

// StatRow signifies a network statistics row in the table
type StatRow struct {
	dest  Server
	stats *ping.Statistics
	err   string
}

func (s StatRow) rtt(datum time.Duration) string {
	ms := datum.Milliseconds()
	color := DefaultGoodColor
	if ms > 200 {
		color = DefaultErrorColor
	} else if ms > 100 && ms < 200 {
		color = DefaultWarnColor
	}
	style := pterm.NewStyle(color)
	return style.Sprintf("%dms", ms)
}

func (s StatRow) loss(datum float64) string {
	color := DefaultGoodColor
	if datum > 10 {
		color = DefaultErrorColor
	} else if datum > 5 && datum < 10 {
		color = DefaultWarnColor
	}
	style := pterm.NewStyle(color, pterm.Italic)
	return style.Sprintf("%0.2f%%", datum)
}

func (s StatRow) name() string {
	var style *pterm.Style
	if s.err != "" {
		style = pterm.NewStyle(pterm.Bold, pterm.FgRed)
	} else {
		style = pterm.NewStyle(pterm.Bold, pterm.FgWhite)
	}
	return style.Sprint(s.dest.Name)
}

func (s StatRow) addr() string {
	var style *pterm.Style
	if s.err != "" {
		style = pterm.NewStyle(pterm.Underscore, pterm.FgRed)
	} else {
		style = pterm.NewStyle(pterm.Underscore)
	}
	return style.Sprint(s.dest.Address)
}

func (s StatRow) error() string {
	style := pterm.NewStyle(pterm.Italic, pterm.FgRed)
	return style.Sprint(s.err)
}

// build adds formatting and styles to the values in the row
func (s StatRow) build() []string {
	timeRecorded := time.Now().Format("2006-01-02 15:04:05")
	if s.err != "" {
		style := pterm.NewStyle(pterm.FgRed)
		return []string{
			s.name(), s.addr(),
			style.Sprint("--"), style.Sprint("--"), style.Sprint("--"), style.Sprint("--"),
			style.Sprint("--"), timeRecorded, s.error(),
		}
	}
	return []string{
		s.name(),
		s.addr(),
		fmt.Sprintf("%d", s.stats.PacketsSent),
		s.loss(s.stats.PacketLoss),
		s.rtt(s.stats.AvgRtt),
		s.rtt(s.stats.MinRtt),
		s.rtt(s.stats.MaxRtt),
		timeRecorded,
		"----",
	}
}

// BuildStatsRow returns formatted network stats data as a row
func BuildStatsRow(dest Server, stats *ping.Statistics, err string) []string {
	row := StatRow{dest: dest, stats: stats, err: err}
	return row.build()
}

// StatsTable signifies the network statistics table in UI
type StatsTable struct {
	// rows is the list of destination stats (including the header)
	rows [][]string
	// addressMap stores the link between the destination address, and it's
	// corresponding row index
	addressMap map[string]int
	// eventChan is used for listening to UI change events
	eventChan chan UIEvent
	// area contains the table
	area *pterm.AreaPrinter
}

func NewStatsTable(destinations []Server, uiEventChan chan UIEvent) *StatsTable {
	if !Config.UIEnabled {
		// Skip if UI is disabled
		return &StatsTable{}
	}
	table := StatsTable{
		rows:       make([][]string, 0),
		addressMap: make(map[string]int),
		eventChan:  uiEventChan,
	}
	// Initialise the area
	var err error
	table.area, err = pterm.DefaultArea.Start()
	if err != nil {
		log.Panic("Failed to draw table area.", err)
	}

	// Add the header row
	table.rows = append(table.rows, StatsTableHeader)
	// Populate the table with the destinations containing empty data
	for idx, dest := range destinations {
		// Create an empty stats row for initialisation
		stats := &ping.Statistics{Addr: dest.Address}
		table.rows = append(table.rows, BuildStatsRow(dest, stats, ""))
		table.addressMap[dest.Address] = idx + 1
	}

	// Draw the table
	tableStr, err := pterm.DefaultTable.WithHasHeader().WithData(table.rows).Srender()
	if err != nil {
		log.Panic("Failed to draw table. ", err)
	}
	table.area.Update(tableStr)

	return &table
}

func (t *StatsTable) listenForChange(ctx context.Context) {
	if !Config.UIEnabled {
		// Skip if UI is disabled
		return
	}
	for {
		select {
		case event := <-t.eventChan:
			// Build the row for the event
			row := BuildStatsRow(event.dest, event.stats, event.err)
			// Update the corresponding row
			t.rows[t.addressMap[event.dest.Address]] = row
			// Re-render the table
			tableStr, _ := pterm.DefaultTable.WithHasHeader().WithData(t.rows).Srender()
			t.area.Update(tableStr)
		case <-ctx.Done():
			Log.Warn("Received termination signal, stopping listening to changes in UI")
			if err := t.area.Stop(); err != nil {
				log.Fatal("Failed to terminate area printer", err)
			}
			return
		}
	}
}

// PrintIntro outputs the program metadata along with title
func PrintIntro() {
	if !Config.UIEnabled {
		return
	}
	pterm.Println()
	if err := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("Ekko ", pterm.NewStyle(pterm.FgLightMagenta)),
		pterm.NewLettersFromStringWithStyle("Network Test", pterm.NewStyle(pterm.FgLightYellow)),
	).Render(); err != nil {
		Log.Panic("Failed to print header", zap.Error(err))
	}
	pterm.Println()
	if Config.Logging.FileEnabled {
		pterm.Info.Println("File logging enabled")
		pterm.Info.Printfln("Results Log path: %s/%s", FileLogDirectory, ResultsLogName)
		pterm.Info.Printfln("Debug Log path: %s/%s", FileLogDirectory, DebugLogName)
	}

	pterm.Println()
}
