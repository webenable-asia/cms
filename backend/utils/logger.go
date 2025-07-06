package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()

	// Set output to stdout
	Logger.SetOutput(os.Stdout)

	// Set log level based on environment
	if os.Getenv("LOG_LEVEL") == "debug" {
		Logger.SetLevel(logrus.DebugLevel)
	} else {
		Logger.SetLevel(logrus.InfoLevel)
	}

	// Use JSON formatter for production
	if os.Getenv("NODE_ENV") == "production" {
		Logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}
}

// LogError logs an error with context
func LogError(err error, message string, fields logrus.Fields) {
	if fields == nil {
		fields = logrus.Fields{}
	}
	fields["error"] = err.Error()
	Logger.WithFields(fields).Error(message)
}

// LogInfo logs an info message with context
func LogInfo(message string, fields logrus.Fields) {
	if fields == nil {
		fields = logrus.Fields{}
	}
	Logger.WithFields(fields).Info(message)
}

// LogWarning logs a warning message with context
func LogWarning(message string, fields logrus.Fields) {
	if fields == nil {
		fields = logrus.Fields{}
	}
	Logger.WithFields(fields).Warn(message)
}

// LogDebug logs a debug message with context
func LogDebug(message string, fields logrus.Fields) {
	if fields == nil {
		fields = logrus.Fields{}
	}
	Logger.WithFields(fields).Debug(message)
}
