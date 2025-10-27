package httpserver

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

type EasyGoHTTPServer struct {
	server *http.Server
	Chi    *chi.Mux
}

func (s *EasyGoHTTPServer) ListenAndServe() error {
	return s.server.ListenAndServe()
}

type NewEasyGoHTTPServerArgs struct {
	Logger *logrus.Logger
	Port   int
}

// customLogFormatter skips logging for health check endpoints
type customLogFormatter struct {
	Logger  *logrus.Logger
	NoColor bool
}

func (l *customLogFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	// Skip logging for health check endpoints
	if r.URL.Path == "/healthz" || r.URL.Path == "/" {
		return &noopLogEntry{}
	}

	// Use the default formatter for other requests
	return &defaultLogEntry{
		Logger:  l.Logger,
		NoColor: l.NoColor,
	}
}

// noopLogEntry does nothing when logging
type noopLogEntry struct{}

func (e *noopLogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	// Do nothing - skip logging
}

func (e *noopLogEntry) Panic(v interface{}, stack []byte) {
	// Do nothing - skip logging
}

// defaultLogEntry provides basic logging functionality
type defaultLogEntry struct {
	Logger  *logrus.Logger
	NoColor bool
}

func (e *defaultLogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	e.Logger.WithFields(logrus.Fields{
		"status":  status,
		"bytes":   bytes,
		"elapsed": elapsed,
	}).Info("HTTP request completed")
}

func (e *defaultLogEntry) Panic(v interface{}, stack []byte) {
	e.Logger.WithFields(logrus.Fields{
		"panic": v,
		"stack": string(stack),
	}).Error("HTTP request panic")
}

func (s *EasyGoHTTPServer) GetHttpServer() *http.Server {
	return s.server
}

func NewEasyGoHTTPServer(args *NewEasyGoHTTPServerArgs) *EasyGoHTTPServer {
	if args.Logger == nil {
		args.Logger = logrus.New()
		args.Logger.SetFormatter(&logrus.JSONFormatter{})
	}

	r := chi.NewRouter()
	// Create a custom logger that skips health check endpoints
	r.Use(middleware.RequestLogger(&customLogFormatter{
		Logger:  args.Logger,
		NoColor: true,
	}))

	server := &http.Server{
		Addr:     fmt.Sprintf(":%d", args.Port),
		Handler:  r,
		ErrorLog: log.New(io.Discard, "", 0), // Disable default logging
	}

	return &EasyGoHTTPServer{
		server: server,
		Chi:    r,
	}
}
