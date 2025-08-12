# Foundation

A reusable Go foundation module that provides common application infrastructure including CLI framework, configuration management, and logging.

## Features

- **CLI Framework**: Built on top of `urfave/cli/v2` with common flags and lifecycle management
- **Configuration Management**: YAML file configuration with environment variable overrides
- **Logging**: Structured logging with configurable levels and formats
- **Signal Handling**: Graceful shutdown handling
- **Environment Variable Support**: Automatic binding and override capabilities

## Installation

```bash
go get github.com/helloworldyuhaiyang/foundation
```

## Quick Start

```go
package main

import (
    "os"
    
    "github.com/helloworldyuhaiyang/foundation/app"
    "github.com/sirupsen/logrus"
    "github.com/urfave/cli/v2"
)

func main() {
    myApp := app.NewApp("my-app", "My application description")
    myApp.SetVersion("1.0.0")
    
    myApp.Init(
        app.WithCommands([]*cli.Command{
            {
                Name: "serve",
                Usage: "Start the server",
                Action: func(c *cli.Context) error {
                    return runServer(c, myApp)
                },
            },
        }),
        // Set environment variable prefix
        app.WithEnvPrefix("APP"), // APP_SERVER_PORT → server.port
    )
    
    if err := myApp.Start(); err != nil {
        logrus.Fatal(err)
    }
}

func runServer(c *cli.Context, myApp *app.App) error {
    // Get configuration - bound from prefixed environment variables!
    port := myApp.Config().GetString("server.port")     // From APP_SERVER_PORT
    dbURL := myApp.Config().GetString("database.url")   // From APP_DATABASE_URL
    
    if port == "" {
        port = "8080"
    }
    
    logrus.Infof("Starting server on port %s", port)
    logrus.Infof("Database URL: %s", dbURL)
    
    // Wait for shutdown signal
    app.WaitForSignal(func(s os.Signal) {
        logrus.Infof("Received signal %v, shutting down gracefully", s)
    })
    
    return nil
}
```

## Configuration

### Configuration File

Create a YAML configuration file (default: `./config/default.yaml`):

```yaml
server:
  port: "8080"
  host: "localhost"

database:
  url: "postgres://localhost/mydb"
  max_connections: 10

log:
  level: "info"
  format: "text"  # or "json"
```

### Environment Variable Overrides

Environment variables automatically override configuration file values using Viper's built-in support:

#### Environment Variable Prefix
Set a prefix for your application's environment variables:

```go
app.WithEnvPrefix("APP") // APP_SERVER_PORT -> server.port
```

- **With prefix**: `APP_SERVER_PORT` → `server.port`
- **With prefix**: `APP_DATABASE_URL` → `database.url`
- **With prefix**: `APP_LOG_LEVEL` → `log.level`

#### Manual Bindings
For custom mappings or environment variables without prefix:

```go
app.WithEnvBindings(map[string]string{
    "server.port": "MY_PORT",           // Custom env var name
    "api.key":     "SECRET_API_KEY",    // Non-standard mapping
})
```

#### Automatic Key Replacement
Viper automatically converts between formats:
- Config file: `server.port`
- Environment: `APP_SERVER_PORT` (with prefix "APP")
- Dots become underscores automatically

### Built-in Flags

The foundation automatically provides these CLI flags:

- `--config, -c`: Configuration file path (default: ./config/default.yaml)
- `--log.level`: Log level (debug, info, warn, error)
- `--log.format`: Log format (text, json)
- `--env`: Environment (dev, test, prod)

## Components

### App

The main application wrapper that provides:

- CLI framework setup
- Configuration loading
- Logger initialization
- Lifecycle management

```go
app := app.NewApp("my-app", "Description")
app.SetVersion("1.0.0")
app.Init(options...)
app.Start()
```

### Config Manager

Configuration management with file and environment support:

```go
config := myApp.Config()
value := config.GetString("key")
config.UnmarshalKey("section", &struct{})
```

### Logger

Structured logging with configurable output:

```go
log := logger.GetLogger("module-name")
log.Info("Message")
log.WithField("key", "value").Error("Error message")
```

## Options

Configure the application using option functions:

- `WithCommands()`: Add CLI commands
- `WithFlags()`: Add custom flags  
- `WithConfigFile()`: Set default config file
- `WithEnvBindings()`: Add environment variable bindings
- `WithContext()`: Set application context
- `AddBefore()`: Add pre-execution hooks
- `AddAfter()`: Add post-execution hooks

## License

MIT License
