# Foundation

A reusable Go foundation module that provides common application infrastructure including CLI framework, configuration management, and logging.

## Features

- **CLI Framework**: Built on top of `urfave/cli/v2` with common flags and lifecycle management
- **Configuration Management**: YAML file configuration with environment variable overrides
- **Enhanced Struct Unmarshaling**: Direct environment variable mapping for struct configuration
- **Logging**: Structured logging with configurable levels and formats
- **Signal Handling**: Graceful shutdown handling
- **Environment Variable Support**: Automatic binding and override capabilities

## Installation

```bash
go get github.com/letusgogo/quick
```

## Quick Start

```go
package main

import (
    "os"
    
    "github.com/letusgogo/quick/app"
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

type ServerConfig struct {
    Port string `mapstructure:"port"`
    Host string `mapstructure:"host"`
}

func runServer(c *cli.Context, myApp *app.App) error {
    // Get configuration - bound from prefixed environment variables!
    dbURL := myApp.Config().GetString("database.url")   // From APP_DATABASE_URL
    
    // Enhanced: Get server config with environment variable sync
    var serverConfig ServerConfig
    envMappings := map[string]string{
        "server.port": "SERVER_PORT",  // Custom env var without prefix
        "server.host": "SERVER_HOST",
    }
    if err := myApp.Config().UnmarshalKeyWithEnv("server", &serverConfig, envMappings); err != nil {
        return err
    }
    
    if serverConfig.Port == "" {
        serverConfig.Port = "8080"
    }
    
    logrus.Infof("Starting server on %s:%s", serverConfig.Host, serverConfig.Port)
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

#### UnmarshalKey with Environment Variable Sync

For struct unmarshaling with environment variable support, use the enhanced method:

```go
type DatabaseConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Username string `mapstructure:"username"`
    Password string `mapstructure:"password"`
}

var dbConfig DatabaseConfig
envMappings := map[string]string{
    "database.host":     "DB_HOST",
    "database.port":     "DB_PORT", 
    "database.username": "DB_USER",
    "database.password": "DB_PASSWORD",
}

// This will automatically sync environment variables before unmarshaling
err := config.UnmarshalKeyWithEnv("database", &dbConfig, envMappings)
```

**Benefits:**
- ✅ No need for pre-binding environment variables
- ✅ Specify mappings only when needed
- ✅ One-line solution for struct + environment variables
- ✅ Explicit control over which variables are mapped

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

// Enhanced: Unmarshal with automatic environment variable sync
type ServerConfig struct {
    Port string `mapstructure:"port"`
    Host string `mapstructure:"host"`
}

var serverConfig ServerConfig
envMappings := map[string]string{
    "server.port": "SERVER_PORT",
    "server.host": "SERVER_HOST",
}
config.UnmarshalKeyWithEnv("server", &serverConfig, envMappings)
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
