// Filename: cmd/api/server.go
package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// serve starts the HTTP server with configured settings.
func (app *applicationDependencies) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),                      // Server address
		Handler:      app.routes(),                                             // HTTP handler routes
		IdleTimeout:  time.Minute,                                              // Idle connection timeout
		ReadTimeout:  5 * time.Second,                                          // Read timeout
		WriteTimeout: 10 * time.Second,                                         // Write timeout
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError), // Error logger
	}

	// Log server startup information.
	app.logger.Info("starting server", "addr", srv.Addr, "env", app.config.environment)
	return srv.ListenAndServe() // Start listening for HTTP requests.
}