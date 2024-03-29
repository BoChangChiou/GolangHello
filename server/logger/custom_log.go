package logger

import (
	"log/slog"
	"os"

	"github.com/natefinch/lumberjack"
	slogmulti "github.com/samber/slog-multi"
)

var (
	logger       *slog.Logger
	defaultLevel = slog.LevelDebug // If don't set to Debug, there is no Debug Level log
)

func init() {
	logFile := &lumberjack.Logger{
		Filename: "C:/Users/User/Documents/golang_test.log",
		MaxSize:  100,
		MaxAge:   1,
	}

	logger = slog.New(
		slogmulti.Fanout(
			slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: defaultLevel}),
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: defaultLevel}),
		),
	)
	slog.SetDefault(logger)
}

func Error(msg string) {
	logger.Error(msg)
}

func Warn(msg string) {
	logger.Warn(msg)
}

func Debug(msg string) {
	logger.Debug(msg)
}

func Info(msg string) {
	logger.Info(msg)
}

func GetLogger() *slog.Logger {
	return logger
}
