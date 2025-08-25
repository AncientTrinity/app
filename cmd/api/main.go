package main

import (
	"flag"
	"log/slog"
	"os"
)

const version = "1.0.0"

// --- Configuration settings ---
type configuration struct {
	port int
	env  string
}

// --- Application struct with DI ---
type application struct {
	config configuration
	logger *slog.Logger
}

func main() {
	// Load config
	cfg := loadConfig()

	// Setup logger
	logger := setupLogger(cfg.env)

	// Initialize app
	app := &application{
		config: cfg,
		logger: logger,
	}

	// Start server
	err := app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

// loadConfig reads configuration from command line flags
func loadConfig() configuration {
	var cfg configuration
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()
	return cfg
}

// setupLogger configures the application logger based on environment
func setupLogger(env string) *slog.Logger {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return logger
}
