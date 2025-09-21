package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"victortillett.net/basic/internal/data"
)

const appVersion = "1.0.0"

type serverConfig struct {
	port           int
	environment    string
	db             struct {
		dsn string
	}
	cors struct {
		trustedOrigins []string
	}
}

type applicationDependencies struct {
	config       serverConfig
	logger       *slog.Logger
	commentModel data.CommentModel
}

func main() {
	var settings serverConfig
	flag.IntVar(&settings.port, "port", 8081, "Server port")
	flag.StringVar(&settings.environment, "env", "development", "Environment")
	flag.StringVar(&settings.db.dsn, "db-dsn", "postgres://user:password@postgres/mydb?sslmode=disable", "PostgreSQL DSN")

	// Pass a space-separated list of origins, e.g. "http://localhost:8080"
	var corsTrustedOrigins string
	flag.StringVar(&corsTrustedOrigins, "cors-trusted-origins", "http://localhost:8080", "Trusted CORS origins (space separated)")
	flag.Parse()

	// Split into slice
	if corsTrustedOrigins != "" {
		settings.cors.trustedOrigins = strings.Fields(corsTrustedOrigins)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(settings)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("database connection pool established")

	app := &applicationDependencies{
		config:       settings,
		logger:       logger,
		commentModel: data.CommentModel{DB: db},
	}

	apiServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", settings.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	shutdownError := make(chan error)

go func() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	s := <-quit

	logger.Info("shutting down server", "signal", s.String())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	shutdownError <- apiServer.Shutdown(ctx)
}()

logger.Info("server started", "addr", apiServer.Addr, "env", settings.environment)

err = apiServer.ListenAndServe()
if !errors.Is(err, http.ErrServerClosed) {
	logger.Error("server error", "error", err)
	os.Exit(1)
}

err = <-shutdownError
if err != nil {
	logger.Error("graceful shutdown failed", "error", err)
	os.Exit(1)
}

logger.Info("server stopped cleanly")

}

func openDB(settings serverConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", settings.db.dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}
