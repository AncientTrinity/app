package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"victortillett.net/basic/internal/data"
)

const appVersion = "1.0.0"

type serverConfig struct {
	port        int
	environment string
	db          struct {
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
	//userModel    data.UserModel
}

func main() {
	var settings serverConfig

	// CLI flags
	flag.IntVar(&settings.port, "port", 8081, "Server port")
	flag.StringVar(&settings.environment, "env", "development", "Environment")
	flag.StringVar(&settings.db.dsn, "db-dsn", "postgres://user:password@postgres/mydb?sslmode=disable", "PostgreSQL DSN")

	var corsTrustedOrigins string
	flag.StringVar(&corsTrustedOrigins, "cors-trusted-origins", "", "Trusted CORS origins (space separated)")

	flag.Parse()

	// ✅ This must run inside main()
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
		//userModel:    data.UserModel{DB: db}, // if you have users.go
	}

	apiServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", settings.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "address", apiServer.Addr, "environment", settings.environment)
	err = apiServer.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
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
