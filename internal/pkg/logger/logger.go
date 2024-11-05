package logger

import (
	"log"
	"os"
)

type Logger struct {
	logger *log.Logger
}

var defaultLogger *Logger

func init() {
	defaultLogger = &Logger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func Info(format string, v ...interface{}) {
	defaultLogger.logger.Printf("[INFO] "+format, v...)
}

func Error(format string, v ...interface{}) {
	defaultLogger.logger.Printf("[ERROR] "+format, v...)
}

func Debug(format string, v ...interface{}) {
	if os.Getenv("DEBUG") == "true" {
		defaultLogger.logger.Printf("[DEBUG] "+format, v...)
	}
}
