package mylog

import (
	"testing"
)

func TestFileLogger(t *testing.T) {
	logger := NewFileLogger(LogLevelDebug, "/Users/zhaofan/Desktop/zz/logs","test")
	logger.Debug("test debug log")
	logger.Warn("test warn log")
	logger.Fatal("test fatal log")
	logger.Error("test error log")
	logger.Info("test info log")
	logger.Close()
}


func TestConsoleLogger(t *testing.T) {
	logger := NewConsoleLogger(LogLevelDebug)
	logger.Debug("test debug log")
	logger.Warn("test warn log")
	logger.Fatal("test fatal log")
	logger.Error("test error log")
	logger.Info("test info log")
}