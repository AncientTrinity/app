package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

 type config struct {
	port int
	dbDSN string
}

 type application struct {
	config config
}

func main() {
	var cfg config

	// Get config from flags
	flag.IntVar(&cfg.port, "port", 8081, "API server port")
	flag.StringVar(&cfg.dbDSN, "db-dsn", "postgres://user:password@postgres/mydb?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	app := &application{
		config: cfg,
}
}
	// Just for debugging now
	fmt.Println("Using database DSN:", cfg.dbDSN)

	// Register routes
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)

	addr := fmt.Sprintf(":%d", cfg.port)
	fmt.Printf("Starting server on %s...\n", addr)

	err := http.ListenAndServe(addr, mux)
	if err != nil {
		fmt.Println("Server error:", err)
		os.Exit(1)
	} 
