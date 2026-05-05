package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/RagnaCron/linko/internal/linkoerr"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type closeFunc func() error

func initializeLogger(logFile string) (*slog.Logger, closeFunc, error) {
	debugHandler := tint.NewHandler(os.Stderr, &tint.Options{
		Level:       slog.LevelDebug,
		ReplaceAttr: linkoerr.ReplaceAttr,
		NoColor:     !isatty.IsTerminal(os.Stderr.Fd()) || !isatty.IsCygwinTerminal(os.Stderr.Fd()),
	})
	if logFile == "" {
		return slog.New(debugHandler), func() error { return nil }, nil
	}

	lumberLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    1,
		MaxAge:     28,
		MaxBackups: 10,
		LocalTime:  false,
		Compress:   true,
	}
	infoHandler := slog.NewJSONHandler(lumberLogger, &slog.HandlerOptions{
		ReplaceAttr: linkoerr.ReplaceAttr,
	})

	return slog.New(slog.NewMultiHandler(debugHandler, infoHandler)), func() error {
		err := lumberLogger.Close()
		if err != nil {
			return fmt.Errorf("could not close lumberLogger: %w", err)
		}

		return nil
	}, nil
}
