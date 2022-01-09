package ui

import (
	"fmt"
	"github.com/go-ping/ping"
	"github.com/pterm/pterm"
	"github.com/soheltarir/ekko/config"
	"time"
)

const (
	DefaultGoodColor  = pterm.FgGreen
	DefaultWarnColor  = pterm.FgLightYellow
	DefaultErrorColor = pterm.FgRed
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

// StatRow signifies a network statistics row in the table
type StatRow struct {
	dest  config.Server
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

// statsRow returns formatted network stats Data as a row
func statsRow(dest config.Server, stats *ping.Statistics, err string) []string {
	row := StatRow{dest: dest, stats: stats, err: err}
	return row.build()
}
