package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	global       *zap.SugaredLogger
	defaultLevel = zap.NewAtomicLevelAt(zap.WarnLevel)
)

func init() {
	SetLogger(New(defaultLevel))
}

func New(lvl zapcore.LevelEnabler, options ...zap.Option) *zap.SugaredLogger {
	if lvl == nil {
		lvl = defaultLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	sink := zapcore.AddSync(os.Stdout)
	return zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			sink,
			lvl,
		),
		options...,
	).Sugar()
}

func SetLogger(l *zap.SugaredLogger) {
	global = l
}

func SetLevel(l zapcore.Level) {
	defaultLevel.SetLevel(l)
}

func Info(ctx context.Context, args ...interface{}) {
	global.Info(args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	global.Warn(args...)
}

func Error(ctx context.Context, args ...interface{}) {
	global.Error(args...)
}

func Errorf(ctx context.Context, template string, args ...interface{}) {
	global.Errorf(template, args...)
}

func Fatal(ctx context.Context, args ...interface{}) {
	global.Fatal(args...)
}
