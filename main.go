package main

import (
    "fmt"
    "net/http"
)

func main() {
    // Health check endpoint
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "OK")
    })

    // Root endpoint
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello from Go Docker server!")
    })

    fmt.Println("Starting server on :8081")
    if err := http.ListenAndServe(":8081", nil); err != nil {
        fmt.Println("Error starting server:", err)
    }
}
