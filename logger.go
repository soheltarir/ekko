package main

import (
	"fmt"
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
	if !Config.Logging.FileEnabled {
		return
	}
	if _, err := os.Stat(l.logDir); os.IsNotExist(err) {
		err := os.Mkdir(l.logDir, FileLogPermission)
		if err != nil {
			log.Panicf("Failed to create logs folder, err: %s", err)
		}
	}
}

// newLogSetup returns a logging setup object which creates a logger with proper configurations
func newLogSetup() logSetup {
	var logDir string
	if Config.Logging.FileLogsDirectory != "" {
		logDir = Config.Logging.FileLogsDirectory
	} else {
		currWd, _ := os.Getwd()
		logDir = fmt.Sprintf("%s/%s", currWd, DefaultFileLogDir)
	}
	setup := logSetup{logDir: logDir, cores: []zapcore.Core{}}
	setup.init()
	return setup
}

func (l logSetup) setupFileCore(logPath string, logLevels []zapcore.Level) zapcore.Core {
	fp, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, FileLogPermission)
	if err != nil {
		log.Panicf("Failed to open/create the specified log file (%s) to write, err: %s", logPath, err)
	}
	encoder := zapcore.NewJSONEncoder(GlobalLogConfig)
	return zapcore.NewCore(encoder, zapcore.AddSync(fp), justLevels(logLevels))
}

func (l *logSetup) addFileDebugCore() {
	if !Config.Logging.FileEnabled {
		return
	}
	path := fmt.Sprintf("%s/%s", l.logDir, DebugLogName)
	l.cores = append(l.cores, l.setupFileCore(path, []zapcore.Level{zap.DebugLevel, zap.WarnLevel}))
}

func (l *logSetup) addFileInfoCore() {
	if !Config.Logging.FileEnabled {
		return
	}
	path := fmt.Sprintf("%s/%s", l.logDir, ResultsLogName)
	l.cores = append(l.cores, l.setupFileCore(path, []zapcore.Level{zap.InfoLevel, zap.ErrorLevel}))
}

func (l *logSetup) addConsoleCore() {
	if !Config.Logging.ConsoleEnabled {
		return
	}
	encoder := zapcore.NewConsoleEncoder(GlobalLogConfig)
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

// FileLogDirectory is a global variable used for displaying the log paths in the UI (if enabled)
var FileLogDirectory string

func init() {
	setup := newLogSetup()
	Log = setup.finish()
	FileLogDirectory = setup.logDir
}
