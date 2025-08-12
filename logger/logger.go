package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

// Config represents logger configuration
type Config struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// DefaultConfig returns default logger configuration
func DefaultConfig() Config {
	return Config{
		Level:  "info",
		Format: "text",
	}
}

// Init initializes the global logger with the given configuration
func Init(config Config) error {
	return InitWithOptions(config, InitOptions{})
}

// InitOptions provides additional options for logger initialization
type InitOptions struct {
	// Output can be set to redirect log output (default: os.Stdout)
	Output *os.File
	// AddTimestamp controls whether to add timestamp to logs (default: true)
	AddTimestamp bool
	// ForceColors controls whether to force colors in text format (default: true for text format)
	ForceColors *bool
	// ReportCaller controls whether to report caller info (default: true)
	ReportCaller bool
}

// InitWithOptions initializes the global logger with configuration and options
func InitWithOptions(config Config, options InitOptions) error {
	// Parse log level
	parsedLevel, err := logrus.ParseLevel(config.Level)
	if err != nil {
		logrus.Warnf("Invalid log level '%s', using info level", config.Level)
		parsedLevel = logrus.InfoLevel
	}
	logrus.SetLevel(parsedLevel)

	// Set output
	output := options.Output
	if output == nil {
		output = os.Stdout
	}
	logrus.SetOutput(output)

	// Set caller reporting
	reportCaller := options.ReportCaller
	logrus.SetReportCaller(reportCaller)

	// Set formatter based on format
	switch config.Format {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				fileName := fmt.Sprintf("%s:%d", path.Base(f.File), f.Line)
				funcName := path.Base(f.Function)
				return funcName, fileName
			},
		})
	case "text":
		forceColors := true
		if options.ForceColors != nil {
			forceColors = *options.ForceColors
		}

		addTimestamp := true
		if !options.AddTimestamp {
			addTimestamp = false
		}

		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: addTimestamp,
			ForceColors:   forceColors,
			PadLevelText:  true,
		})
	default:
		return fmt.Errorf("unsupported log format: %s", config.Format)
	}

	logrus.Infof("Logger initialized with level=%s, format=%s", config.Level, config.Format)
	return nil
}

// NewLogger creates a new logger instance with the given module name
func NewLogger(module string) *logrus.Entry {
	return logrus.WithFields(map[string]interface{}{
		"module": module,
	})
}

// GetLogger returns a logger with the given module name
func GetLogger(module string) *logrus.Entry {
	return NewLogger(module)
}
