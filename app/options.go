package app

import (
	"context"

	"github.com/urfave/cli/v2"
)

// Options for App
type Options struct {
	// Configuration file path
	ConfigFile string

	// Environment variable prefix (e.g., "APP" for APP_SERVER_PORT)
	EnvPrefix string

	// Command line flags
	Flags []cli.Flag

	// Sub commands
	Commands []*cli.Command

	// Before and After functions
	Before []func(*cli.Context) error
	After  []func(*cli.Context) error

	// Context for the application
	Context context.Context

	// Environment variable bindings for configuration
	EnvBindings map[string]string
}

// NewOptions creates a new Options instance with default values
func NewOptions() *Options {
	return &Options{
		ConfigFile:  "",
		EnvPrefix:   "",
		Flags:       nil,
		Commands:    nil,
		Before:      nil,
		After:       nil,
		Context:     context.Background(),
		EnvBindings: make(map[string]string),
	}
}

// Option is a function that configures Options
type Option func(*Options)

// WithConfigFile sets the configuration file path
func WithConfigFile(configFile string) Option {
	return func(o *Options) {
		o.ConfigFile = configFile
	}
}

// WithEnvPrefix sets the environment variable prefix
// Example: WithEnvPrefix("APP") means APP_SERVER_PORT maps to server.port
func WithEnvPrefix(prefix string) Option {
	return func(o *Options) {
		o.EnvPrefix = prefix
	}
}

// WithCommands sets the CLI commands
func WithCommands(commands []*cli.Command) Option {
	return func(o *Options) {
		o.Commands = commands
	}
}

// WithFlags sets the CLI flags
func WithFlags(flags []cli.Flag) Option {
	return func(o *Options) {
		o.Flags = flags
	}
}

// WithContext sets the application context
func WithContext(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// WithEnvBindings sets environment variable bindings
func WithEnvBindings(bindings map[string]string) Option {
	return func(o *Options) {
		o.EnvBindings = bindings
	}
}

// AddBefore adds a before function
func AddBefore(before func(*cli.Context) error) Option {
	return func(o *Options) {
		o.Before = append(o.Before, before)
	}
}

// AddAfter adds an after function
func AddAfter(after func(*cli.Context) error) Option {
	return func(o *Options) {
		o.After = append(o.After, after)
	}
}

// AddEnvBinding adds a single environment variable binding
func AddEnvBinding(key, envVar string) Option {
	return func(o *Options) {
		if o.EnvBindings == nil {
			o.EnvBindings = make(map[string]string)
		}
		o.EnvBindings[key] = envVar
	}
}
