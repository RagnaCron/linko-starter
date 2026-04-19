package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
)

type closeFunc func() error

func initializeLogger(logFile string) (*slog.Logger, closeFunc, error) {
	debugHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	if logFile == "" {
		return slog.New(debugHandler), func() error { return nil }, nil
	}

	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return nil, func() error { return nil }, fmt.Errorf("failed to open log file: %w", err)
	}
	bufferedFile := bufio.NewWriterSize(file, 8192)
	infoHandler := slog.NewJSONHandler(bufferedFile, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	return slog.New(slog.NewMultiHandler(debugHandler, infoHandler)), func() error {
		err := bufferedFile.Flush()
		if err != nil {
			return fmt.Errorf("could not flush buffer to file: %w", err)
		}
		err = file.Close()
		if err != nil {
			return fmt.Errorf("could not close log file: %w", err)
		}
		return nil
	}, nil
}
