package main

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

// BadLog - неэффективный логгер что-то хуже придумать не получилось.

type BadLog struct {
	wr *bufio.Writer
}

func NewBadLog(w io.Writer) *BadLog {
	return &BadLog{
		wr: bufio.NewWriter(w),
	}
}
func (l *BadLog) Info(msg any) {
	s := fmt.Sprintf("%v INFO:%v", time.Now().Format(time.DateTime), msg.(string))
	fmt.Fprintln(l.wr, s)
}

// На удивление logrus самый не эффективный по бэнчмаркам.
// не совсем честное решение, но...
type VeryBadLog struct {
	logger *logrus.Logger
}

func NewVeryBadLog(w io.Writer) *VeryBadLog {
	l := logrus.New()
	l.Out = w
	return &VeryBadLog{
		logger: l,
	}
}
func (l *VeryBadLog) Info(msg any) {
	s := fmt.Sprintf("%v INFO:%v", time.Now().Format(time.DateTime), msg.(string))
	l.logger.Info(s)

}

// GoodLog - эффективный логгер
type GoodLog struct {
	wr  io.Writer
	buf []byte
}

func NewGoodLog(w io.Writer) *GoodLog {
	return &GoodLog{
		wr:  w,
		buf: make([]byte, 0, 256),
	}
}
func (l *GoodLog) Info(msg string) {
	l.buf = append(l.buf, time.Now().AppendFormat(l.buf, time.RFC3339)...)
	l.buf = append(l.buf, ' ')
	l.buf = append(l.buf, []byte(msg)...)
	l.wr.Write(l.buf)
	l.buf = l.buf[:0]
}

func main() {

	logger1 := slog.New(slog.NewTextHandler(os.Stderr, nil))
	logger2 := logrus.New()
	logger3 := zerolog.New(os.Stderr)
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stderr"}
	logger4, _ := cfg.Build()
	defer logger4.Sync()
	logger5 := NewGoodLog(os.Stderr)
	logger1.Info("happy", "lama", "1")
	logger2.Info("happy", "lama", "2")
	logger3.Info().Msg("happy lama 3")
	logger4.Info("happy lama 4")
	logger5.Info("happy lama 5")

}
