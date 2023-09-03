package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/xid"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nkmr-jp/zl"
)

// handler receives requests and performs processing.
func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now() // Start measuring processing time

	traceID := xid.New().String()             // ID for tracing request scope logs
	ctx := SetNewLogger(r.Context(), traceID) // Create and set the context scope logger
	cl := GetLogger(ctx)                      // Get logger from context

	cl.Info("REQUEST_RECEIVED")

	// Defer function to recover from panic
	defer func(cl *zl.Logger) {
		if err := recover(); err != nil {
			cl.Err("PANIC_RECOVERED", fmt.Errorf("recovered from panic: %v", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}(cl)

	doSomething(ctx)
	doSomething2(ctx)

	_, err := fmt.Fprintf(w, "request trace_id: %s", traceID)
	if err != nil {
		cl.Err("WRITE_ERROR", err)
	}
	cl.Info("REQUEST_ALL_PROCESS_COMPLETE", DurationField(start)) // Recording processing time
}

// Do something within the context of the request
func doSomething(ctx context.Context) {
	start := time.Now()  // Start measuring processing time
	cl := GetLogger(ctx) // Get logger from context

	time.Sleep(time.Second)

	cl.Debug("DO_SOMETHING1", DurationField(start)) // Recording processing time
}

func doSomething2(ctx context.Context) {
	start := time.Now()  // Start measuring processing time
	cl := GetLogger(ctx) // Get logger from context

	time.Sleep(time.Millisecond * 500)

	cl.Debug("DO_SOMETHING2", DurationField(start)) // Recording processing time
}

// initGlobalLogger initializes the global scope logger.
func initGlobalLogger() {
	// Set log level and output destination according to the environment variable
	switch os.Getenv("ENV") {
	case "production":
		zl.SetLevel(zl.InfoLevel)
		zl.SetOutput(zl.ConsoleOutput) // Output to console in json format
	case "development":
		zl.SetLevel(zl.DebugLevel)
		zl.SetOutput(zl.ConsoleOutput) // Output to console in json format
	case "local":
		zl.SetLevel(zl.DebugLevel)
		zl.SetOutput(zl.PrettyOutput) // Output to file and console in colored simple log format
	default:
		log.Fatal("ENV is not set")
	}
	zl.SetOmitKeys(zl.HostnameKey)
	zl.SetConsoleFields(TraceIDFieldKey, DurationFieldKey) // Set fields to be output to console
	zl.Init()                                              // Initialize global logger
}

// main starts the server.
// When it receives a signal, it shuts down the server.
func main() {
	initGlobalLogger()
	defer zl.Sync() // flush log buffer

	// Create a channel to receive signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Start HTTP server with goroutine
	http.HandleFunc("/", handler)
	srv := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: http.DefaultServeMux,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zl.FatalErr("LISTEN_AND_SERVE_FAILED", err)
		}
	}()

	// Wait until signal is received
	sig := <-signalChan
	zl.Info("SIGNAL_CAUGHT", zl.Consolef("signal: %s", sig.String()))

	// Create context with timeout for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gracefully shut down the server. If it does not finish within 10 seconds, an error will be returned.
	if err := srv.Shutdown(ctx); err != nil {
		zl.FatalErr("SERVER_SHUTDOWN_FAILED", err)
	}
	zl.Info("SERVER_SHUTDOWN_SUCCESS")
}
