package logger

import (
	"fmt"
	"github.com/soheltarir/ekko/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

// justLevel returns a LevelEnabler that only enables the specified levels
func justLevels(levels []zapcore.Level) zap.LevelEnablerFunc {
	return func(level zapcore.Level) bool {
		for _, l := range levels {
			if l == level {
				return true
			}
		}
		return false
	}
}

// logSetup contains all the configuration for creating the global logger for the application
type logSetup struct {
	cores []zapcore.Core
	// logDir is the root folder wherein the file logs would reside
	logDir string
}

// init initialises some attributes of the setup object
func (l logSetup) init() {
	// Create the file logs directory if it doesn't exist
	if !config.Config.Logging.FileEnabled {
		return
	}
	if _, err := os.Stat(l.logDir); os.IsNotExist(err) {
		err := os.Mkdir(l.logDir, fileLogPermission)
		if err != nil {
			log.Panicf("Failed to create logs folder, err: %s", err)
		}
	}
}

// newLogSetup returns a logging setup object which creates a logger with proper configurations
func newLogSetup(fileLogDir string) logSetup {
	setup := logSetup{logDir: fileLogDir, cores: []zapcore.Core{}}
	setup.init()
	return setup
}

func (l logSetup) setupFileCore(logPath string, logLevels []zapcore.Level) zapcore.Core {
	fp, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, fileLogPermission)
	if err != nil {
		log.Panicf("Failed to open/create the specified log file (%s) to write, err: %s", logPath, err)
	}
	encoder := zapcore.NewJSONEncoder(Config)
	return zapcore.NewCore(encoder, zapcore.AddSync(fp), justLevels(logLevels))
}

func (l *logSetup) addFileDebugCore() {
	if !config.Config.Logging.FileEnabled {
		return
	}
	path := fmt.Sprintf("%s/%s", l.logDir, debugLogName)
	LogPath.Debug = path
	l.cores = append(l.cores, l.setupFileCore(path, []zapcore.Level{zap.DebugLevel, zap.WarnLevel}))
}

func (l *logSetup) addFileInfoCore() {
	if !config.Config.Logging.FileEnabled {
		return
	}
	path := fmt.Sprintf("%s/%s", l.logDir, resultsLogName)
	LogPath.Results = path
	l.cores = append(l.cores, l.setupFileCore(path, []zapcore.Level{zap.InfoLevel, zap.ErrorLevel}))
}

func (l *logSetup) addConsoleCore() {
	if !config.Config.Logging.ConsoleEnabled {
		return
	}
	encoder := zapcore.NewConsoleEncoder(Config)
	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zap.DebugLevel)
	l.cores = append(l.cores, core)
}

func (l logSetup) finish() *zap.Logger {
	l.addFileInfoCore()
	l.addFileDebugCore()
	l.addConsoleCore()
	core := zapcore.NewTee(l.cores...)
	return zap.New(core).WithOptions(zap.OnFatal(zapcore.WriteThenNoop))
}

var Log *zap.Logger

// logPath contains the file location paths for different kind of logs being used
type logPath struct {
	Results string
	Debug   string
}

// LogPath is the global variable which stores log paths
var LogPath = new(logPath)

func init() {
	setup := newLogSetup(setupFileLogDirectory())
	Log = setup.finish()
}
