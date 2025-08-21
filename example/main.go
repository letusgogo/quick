package main

import (
	"os"
	"time"

	"github.com/letusgogo/quick/app"
	"github.com/letusgogo/quick/logger"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

func main() {
	// Create a new application
	myApp := app.NewApp("example-app", "Example application using foundation module")
	myApp.SetVersion("1.0.0")

	// Initialize with commands and configuration
	myApp.Init(
		// Add commands
		app.WithCommands([]*cli.Command{
			{
				Name:  "serve",
				Usage: "Start the HTTP server",
				Action: func(c *cli.Context) error {
					return runServer(c, myApp)
				},
			},
			{
				Name:  "worker",
				Usage: "Start background worker",
				Action: func(c *cli.Context) error {
					return runWorker(c, myApp)
				},
			},
		}),

		// Add custom flags
		app.WithFlags([]cli.Flag{
			&cli.StringFlag{
				Name:  "mode",
				Value: "production",
				Usage: "application mode",
			},
		}),

		// Set environment variable prefix
		// With prefix "APP": APP_SERVER_PORT -> server.port, APP_DATABASE_URL -> database.url
		app.WithEnvPrefix("APP"),

		// Manual bindings for non-standard environment variable names
		app.WithEnvBindings(map[string]string{
			"custom.api.key": "MY_CUSTOM_API_KEY", // Custom mapping example
		}),

		// Add before hooks
		app.AddBefore(func(c *cli.Context) error {
			logger.GetLogger("main").Infof("Application starting in %s mode", c.String("mode"))
			return nil
		}),

		// Add after hooks
		app.AddAfter(func(c *cli.Context) error {
			logger.GetLogger("main").Info("Application finished")
			return nil
		}),
	)

	// Start the application
	if err := myApp.Start(); err != nil {
		logrus.Fatal(err)
	}
}

func runServer(c *cli.Context, myApp *app.App) error {
	log := logger.GetLogger("server")
	// Get configuration values from env start with APP_
	databaseURL := myApp.Config().GetString("database.url")
	log.Infof("Database URL: %s", databaseURL)

	// Get configuration values
	port := myApp.Config().GetString("server.port")
	host := myApp.Config().GetString("server.host")
	mode := c.String("mode")

	// Get configuration values from struct
	var serverConfig ServerConfig
	envMappings := map[string]string{
		"server.port": "SERVER_PORT",
		"server.host": "SERVER_HOST",
	}
	if err := myApp.Config().UnmarshalKeyWithEnv("server", &serverConfig, envMappings); err != nil {
		log.Fatalf("Failed to unmarshal server config: %v", err)
	}
	log.Infof("Server config from struct, %s:%s", serverConfig.Host, serverConfig.Port)

	log.Infof("Starting HTTP server on %s:%s in %s mode", host, port, mode)

	// Simulate server startup
	log.Info("Server started successfully")

	// Log some configuration for debugging - these are bound from environment variables!
	myApp.Config().LogConfigValue("server.port")        // From APP_SERVER_PORT
	myApp.Config().LogConfigValue("server.host")        // From APP_SERVER_HOST
	myApp.Config().LogConfigValue("database.url")       // From APP_DATABASE_URL
	myApp.Config().LogConfigValue("redis.addr")         // From APP_REDIS_ADDR
	myApp.Config().LogConfigValue("worker.concurrency") // From APP_WORKER_CONCURRENCY

	// Wait for shutdown signal
	app.WaitForSignal(func(s os.Signal) {
		log.Infof("Received signal %v, shutting down HTTP server gracefully", s)
		// Here you would normally stop your HTTP server
		time.Sleep(1 * time.Second) // Simulate graceful shutdown
		log.Info("HTTP server stopped")
	})

	return nil
}

func runWorker(c *cli.Context, myApp *app.App) error {
	log := logger.GetLogger("worker")

	// Get configuration values
	concurrency := myApp.Config().GetInt("worker.concurrency")
	if concurrency == 0 {
		concurrency = 5
	}

	mode := c.String("mode")

	log.Infof("Starting background worker with concurrency=%d in %s mode", concurrency, mode)

	// Simulate worker tasks
	for i := 0; i < concurrency; i++ {
		go func(id int) {
			workerLog := log.WithField("worker_id", id)
			workerLog.Info("Worker started")

			// Simulate work
			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					workerLog.Info("Processing task...")
				}
			}
		}(i)
	}

	log.Info("All workers started successfully")

	// Wait for shutdown signal
	app.WaitForSignal(func(s os.Signal) {
		log.Infof("Received signal %v, shutting down workers gracefully", s)
		// Here you would normally stop your workers
		time.Sleep(2 * time.Second) // Simulate graceful shutdown
		log.Info("All workers stopped")
	})

	return nil
}
