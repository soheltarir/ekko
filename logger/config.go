package logger

import "go.uber.org/zap/zapcore"

var Config = zapcore.EncoderConfig{
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
