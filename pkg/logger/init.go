package logger

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"log"
	"os"
	"time"
)

// NewLoggerWithFile returns a new logger with file output.
//
// Sample usage: NewLoggerWithFile("/var/log/myproject/myproject.log")
func NewLoggerWithFile(filepath string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		filepath,
		"stderr", // for also writing to terminal
	}
	return cfg.Build()
}

// CreateNewFile checks filepath if it exists or not, if it does not exist, creates it.
//
// Use it like this: CreateNewFile("./logs/backend")
//
// This will create a backend_(today'date).log file.
//
// Returns filepath in the end, such as ./logs/backend_20210101.log
func CreateNewFile(filepath string) string {
	LogFilepath := fmt.Sprintf("%s_%s.log", filepath, time.Now().Format("20060102"))
	if _, err := os.Stat(LogFilepath); errors.Is(err, os.ErrNotExist) {
		_, err = os.Create(LogFilepath)
		if err != nil {
			log.Println(err)
		}
	}
	return LogFilepath
}
