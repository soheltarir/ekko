package logger

import (
	"fmt"
	"github.com/soheltarir/ekko/config"
	"os"
)

const (
	debugLogName      = "debug.ndjson"
	resultsLogName    = "results.ndjson"
	fileLogPermission = 0644
)

func setupFileLogDirectory() string {
	if config.Config.Logging.FileLogsDirectory != "" {
		return config.Config.Logging.FileLogsDirectory
	} else {
		currWd, _ := os.Getwd()
		return fmt.Sprintf("%s/logs", currWd)
	}
}
