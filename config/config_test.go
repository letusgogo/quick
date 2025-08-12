package config

import (
	"os"
	"testing"
)

func TestEnvironmentVariableOverrides(t *testing.T) {
	// Set test environment variables with prefix
	os.Setenv("TEST_SERVER_PORT", "9090")
	os.Setenv("TEST_DATABASE_URL", "postgres://test:5432/testdb")
	defer os.Unsetenv("TEST_SERVER_PORT")
	defer os.Unsetenv("TEST_DATABASE_URL")

	manager := NewManager()
	manager.SetEnvPrefix("TEST")
	manager.SetupEnvironmentOverrides()

	// Test with prefix: TEST_SERVER_PORT -> server.port
	port := manager.GetString("server.port")
	if port != "9090" {
		t.Errorf("Expected server.port to be '9090', got '%s'", port)
	}

	// Test with prefix: TEST_DATABASE_URL -> database.url
	dbURL := manager.GetString("database.url")
	if dbURL != "postgres://test:5432/testdb" {
		t.Errorf("Expected database.url to be 'postgres://test:5432/testdb', got '%s'", dbURL)
	}
}

func TestEnvironmentVariableWithoutPrefix(t *testing.T) {
	// Set test environment variables without prefix
	os.Setenv("SERVER_PORT", "8080")
	defer os.Unsetenv("SERVER_PORT")

	manager := NewManager()
	manager.SetupEnvironmentOverrides()

	// Without prefix, need manual binding
	manager.BindEnv("server.port", "SERVER_PORT")

	port := manager.GetString("server.port")
	if port != "8080" {
		t.Errorf("Expected server.port to be '8080', got '%s'", port)
	}
}

func TestManualEnvironmentBinding(t *testing.T) {
	// Set test environment variable
	os.Setenv("MY_CUSTOM_PORT", "3000")
	defer os.Unsetenv("MY_CUSTOM_PORT")

	manager := NewManager()
	manager.SetupEnvironmentOverrides()

	// Manual binding for custom environment variable name
	bindings := map[string]string{
		"server.port": "MY_CUSTOM_PORT",
	}
	manager.BindEnvs(bindings)

	port := manager.GetString("server.port")
	if port != "3000" {
		t.Errorf("Expected server.port to be '3000', got '%s'", port)
	}
}
