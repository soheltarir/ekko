package main

import (
	"github.com/pterm/pterm"
	"go.uber.org/zap/zapcore"
)

const (
	DefaultColor      = pterm.FgWhite
	DefaultGoodColor  = pterm.FgGreen
	DefaultWarnColor  = pterm.FgLightYellow
	DefaultErrorColor = pterm.FgRed
	DefaultFileLogDir = "logs"
	ResultsLogName    = "results.ndjson"
	DebugLogName      = "debug.ndjson"
	FileLogPermission = 0644
)

type ConsumerStatus int

const (
	ConsumerNotStarted ConsumerStatus = iota
	ConsumerRunning
	ConsumerStopped
)

var GlobalLogConfig = zapcore.EncoderConfig{
	MessageKey:     "message",
	LevelKey:       "severity",
	TimeKey:        "timestamp",
	NameKey:        "logger",
	CallerKey:      zapcore.OmitKey,
	FunctionKey:    zapcore.OmitKey,
	StacktraceKey:  zapcore.OmitKey,
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.LowercaseLevelEncoder,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.MillisDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}
