package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

// setupLogger initialises the logger to be used in the application
func setupLogger() *zap.Logger {
	var enabledCores []zapcore.Core
	if Config.Logging.FileEnabled {
		fp, err := os.OpenFile(Config.Logging.FileOutput, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Panicf("Failed to open/create the specified log file (%s) to write, err: %s",
				Config.Logging.FileOutput, err)
		}
		fileEncoder := zapcore.NewJSONEncoder(GlobalLogConfig)
		fileCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(fp), zap.InfoLevel)
		enabledCores = append(enabledCores, fileCore)
	}
	if Config.Logging.ConsoleEnabled {
		consoleEncoder := zapcore.NewConsoleEncoder(GlobalLogConfig)
		consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel)
		enabledCores = append(enabledCores, consoleCore)
	}
	core := zapcore.NewTee(enabledCores...)
	return zap.New(core).WithOptions(zap.OnFatal(zapcore.WriteThenNoop))
}

var Log *zap.Logger

func init() {
	Log = setupLogger()
	Log.Debug("Logger initialised.")
}
