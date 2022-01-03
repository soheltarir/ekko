package ui

import (
	"github.com/pterm/pterm"
	"github.com/soheltarir/ekko/config"
	"github.com/soheltarir/ekko/logger"
)

func header() string {
	header, err := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("Ekko ", pterm.NewStyle(pterm.FgLightMagenta)),
		pterm.NewLettersFromStringWithStyle("Network Test", pterm.NewStyle(pterm.FgLightYellow)),
	).Srender()
	if err != nil {

	}
	return header
}

func programInfo(status config.ConsumerStatus) []interface{} {
	var lines []interface{}
	if config.Config.Logging.FileEnabled {
		lines = append(lines,
			pterm.Info.Sprintln("File logging enabled"),
			pterm.Info.Sprintfln("Results Log path: %s", logger.LogPath.Results),
			pterm.Info.Sprintfln("Debug Log path: %s", logger.LogPath.Debug),
		)
	}
	lines = append(lines, pterm.Info.Sprintln(status))
	return lines
}

func networkStatsTable(rows [][]string) (string, error) {
	if len(rows) == 0 {
		return "", nil
	}
	return pterm.DefaultTable.WithHasHeader(true).WithBoxed(true).
		WithData(rows).Srender()
}
