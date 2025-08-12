package app

import (
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/letusgogo/quick/config"
	"github.com/letusgogo/quick/logger"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// App represents the application
type App struct {
	Name    string
	Usage   string
	Version string
	log     *logrus.Entry
	opt     *Options
	app     *cli.App
	config  *config.Manager
}

// NewApp creates a new application instance
func NewApp(name, usage string) *App {
	return &App{
		Name:   name,
		Usage:  usage,
		config: config.NewManager(),
		log:    logger.GetLogger(name),
	}
}

// SetVersion sets the application version
func (a *App) SetVersion(version string) {
	a.Version = version
}

// Init initializes the application with the given options
func (a *App) Init(opts ...Option) {
	a.app = cli.NewApp()
	a.app.EnableBashCompletion = true
	a.app.Name = a.Name
	a.app.Usage = a.Usage
	if a.Version != "" {
		a.app.Version = a.Version
	}

	a.opt = NewOptions()
	for _, opt := range opts {
		opt(a.opt)
	}

	a.app.Commands = a.opt.Commands
	a.app.Flags = a.opt.Flags

	// Add built-in flags
	a.addBuiltinFlags()

	// Set up before and after handlers
	a.setupHandlers()
}

// addBuiltinFlags adds common flags that most applications need
func (a *App) addBuiltinFlags() {
	defaultConfig := "./config/default.yaml"
	if a.opt.ConfigFile != "" {
		defaultConfig = a.opt.ConfigFile
	}

	builtinFlags := []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			Aliases:     []string{"c"},
			Value:       defaultConfig,
			DefaultText: defaultConfig,
			Usage:       "config file path",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        "log.level",
			Value:       "info",
			DefaultText: "info",
			Usage:       "log level (debug, info, warn, error)",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        "log.format",
			Value:       "text",
			DefaultText: "text",
			Usage:       "log format (text, json)",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        "env",
			Value:       "dev",
			DefaultText: "dev",
			Usage:       "environment (dev, test, prod)",
			Required:    false,
		},
	}

	a.app.Flags = append(a.app.Flags, builtinFlags...)
}

// setupHandlers sets up before and after handlers
func (a *App) setupHandlers() {
	a.app.Before = func(c *cli.Context) error {
		// Initialize configuration
		if err := a.initConfig(c); err != nil {
			return err
		}

		// Initialize logger
		if err := a.initLogger(c); err != nil {
			return err
		}

		// Run user-defined before functions
		for _, before := range a.opt.Before {
			if err := before(c); err != nil {
				return err
			}
		}

		return nil
	}

	a.app.After = func(c *cli.Context) error {
		// Run user-defined after functions
		for _, after := range a.opt.After {
			if err := after(c); err != nil {
				return err
			}
		}
		return nil
	}
}

// initConfig initializes configuration management
func (a *App) initConfig(c *cli.Context) error {
	// Setup environment variable overrides using Viper's built-in support
	a.config.SetupEnvironmentOverrides()

	// Set environment prefix if specified in options
	if a.opt.EnvPrefix != "" {
		a.config.SetEnvPrefix(a.opt.EnvPrefix)
	}

	// Load configuration file first
	configFile := c.String("config")
	if err := a.config.LoadFromFile(configFile); err != nil {
		// Not a fatal error, we can continue with environment variables
		a.log.Warnf("Failed to load config file: %v", err)
	}

	// Bind user-defined environment variables for specific mappings
	if len(a.opt.EnvBindings) > 0 {
		a.config.BindEnvs(a.opt.EnvBindings)
	}

	// Bind common environment variables that don't follow the standard pattern
	commonBindings := map[string]string{
		"log.level":  "LOG_LEVEL",
		"log.format": "LOG_FORMAT",
		"env":        "ENV",
	}
	a.config.BindEnvs(commonBindings)

	return nil
}

// initLogger initializes the logger
func (a *App) initLogger(c *cli.Context) error {
	// Get log configuration from CLI flags or config file
	logLevel := c.String("log.level")
	if logLevel == "" {
		logLevel = a.config.GetString("log.level")
	}
	if logLevel == "" {
		logLevel = "info"
	}

	logFormat := c.String("log.format")
	if logFormat == "" {
		logFormat = a.config.GetString("log.format")
	}
	if logFormat == "" {
		logFormat = "text"
	}

	// Initialize logger
	loggerConfig := logger.Config{
		Level:  logLevel,
		Format: logFormat,
	}

	options := logger.InitOptions{
		ReportCaller: true,
		AddTimestamp: true,
	}

	return logger.InitWithOptions(loggerConfig, options)
}

// Start starts the application
func (a *App) Start() error {
	if a.app == nil {
		panic("please call Init() first")
	}

	err := a.app.Run(os.Args)
	if err != nil {
		a.log.Fatal(err)
		return err
	}

	return nil
}

// Config returns the configuration manager
func (a *App) Config() *config.Manager {
	if a.config == nil {
		panic("configuration not initialized, call Init() first")
	}
	return a.config
}

// WaitForSignal waits for termination signals and calls the provided function
func WaitForSignal(stopFunc func(os.Signal)) {
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	defer func() {
		if e := recover(); e != nil {
			logrus.Errorf("crashed, err: %s stack:%s", e, string(debug.Stack()))
		}
	}()

	recvSignal := <-signalChan
	logrus.Infof("received signal: %v", recvSignal)
	stopFunc(recvSignal)
}
