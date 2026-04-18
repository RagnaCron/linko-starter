package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

type closeFunc func() error

func initializeLogger(logFile string) (*log.Logger, closeFunc, error) {
	if logFile == "" {
		return log.New(os.Stderr, "", log.LstdFlags), func() error { return nil }, nil
	}

	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return nil, func() error { return nil }, fmt.Errorf("failed to open log file: %w", err)
	}
	bufferedFile := bufio.NewWriterSize(file, 8192)
	multiWriter := io.MultiWriter(os.Stderr, bufferedFile)
	return log.New(multiWriter, "", log.LstdFlags), func() error {
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
