package main

import (
	"io"
	"log/slog"
	"testing"

	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func BenchmarkSlog(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark test1")
	}
}

func BenchmarkLogrus(b *testing.B) {
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark test2")
	}
}

func BenchmarkZero(b *testing.B) {
	logger := zerolog.New(io.Discard)
	for i := 0; i < b.N; i++ {
		logger.Info().Msg("benchmark test3")
	}
}

func BenchmarkZap(b *testing.B) {
	cfg := zap.NewProductionConfig()
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg.EncoderConfig),
		zapcore.AddSync(io.Discard),
		zapcore.InfoLevel,
	)
	logger := zap.New(core)
	defer logger.Sync()

	for i := 0; i < b.N; i++ {
		logger.Info("benchmark test4")
	}
}
func BenchmarkGoodlog(b *testing.B) {
	logger := NewGoodLog(io.Discard)
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark test5")
	}
}

func BenchmarkBadlog(b *testing.B) {
	logger := NewBadLog(io.Discard)
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark test6")
	}
}

func BenchmarkVeryBadlog(b *testing.B) {
	logger := NewVeryBadLog(io.Discard)
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark test6")
	}
}
