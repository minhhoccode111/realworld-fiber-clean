// Package httpserver implements HTTP server.
package httpserver

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/logger"
	"golang.org/x/sync/errgroup"
)

const (
	_defaultAddr            = ":80"
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultShutdownTimeout = 3 * time.Second
)

// Server -.
type Server struct {
	ctx context.Context
	eg  *errgroup.Group

	App    *gin.Engine
	notify chan error
	server *http.Server

	address         string
	readTimeout     time.Duration
	writeTimeout    time.Duration
	shutdownTimeout time.Duration

	logger logger.Interface
}

// New -.
func New(l logger.Interface, opts ...Option) *Server {
	group, ctx := errgroup.WithContext(context.Background())
	group.SetLimit(1) // Run only one goroutine

	s := &Server{
		ctx:             ctx,
		eg:              group,
		App:             nil,
		notify:          make(chan error, 1),
		address:         _defaultAddr,
		readTimeout:     _defaultReadTimeout,
		writeTimeout:    _defaultWriteTimeout,
		shutdownTimeout: _defaultShutdownTimeout,
		logger:          l,
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	s.App = gin.New() // Initialize Gin engine

	s.server = &http.Server{
		Addr:         s.address,
		Handler:      s.App, // Assign Gin engine as handler
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
	}

	return s
}

// Start -.
func (s *Server) Start() {
	s.eg.Go(func() error {
		err := s.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.notify <- err
			close(s.notify)
			return err
		}
		return nil
	})

	s.logger.Info("http server - Server - Started")
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() error {
	var shutdownErrors []error

	// Create a context with the shutdown timeout
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil && !errors.Is(err, context.Canceled) {
		s.logger.Error(err, "http server - Server - Shutdown - s.server.Shutdown")
		shutdownErrors = append(shutdownErrors, err)
	}

	// Wait for all goroutines to finish and get any error
	err = s.eg.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		s.logger.Error(err, "http server - Server - Shutdown - s.eg.Wait")
		shutdownErrors = append(shutdownErrors, err)
	}

	s.logger.Info("http server - Server - Shutdown")

	return errors.Join(shutdownErrors...)
}
