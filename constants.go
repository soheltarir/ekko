package main

import (
	"github.com/pterm/pterm"
	"go.uber.org/zap/zapcore"
)

const (
	DefaultFileLogPath = "logs/results.ndjson"
	DefaultColor       = pterm.FgWhite
	DefaultWarnColor   = pterm.FgLightYellow
	DefaultErrorColor  = pterm.FgRed
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
