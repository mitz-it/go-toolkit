package hosting

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/mitz-it/go-toolkit/logger"
)

type Host struct {
	Router          http.Handler
	Workers         []Worker
	ShutdownTimeout time.Duration
	Addr            string
	server          *http.Server
	wg              *sync.WaitGroup
}

func (h *Host) Start() error {
	if h.ShutdownTimeout == 0 {
		h.ShutdownTimeout = 30 * time.Second
	}

	switch true {
	case h.Router != nil:
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		h.wg = &sync.WaitGroup{}

		h.work(ctx)

		h.serve(ctx, cancel)

		<-ctx.Done()
		logger.Info(context.Background()).Msg("host shutdown started")
		defer logger.Info(context.Background()).Msg("host shutdown completed")
		h.wg.Wait()

		ctx, cancel = context.WithTimeout(context.Background(), h.ShutdownTimeout)
		defer cancel()

		if err := h.server.Shutdown(ctx); err != nil {
			logger.Info(context.Background()).Msg("HTTP server was forced to shutdown")
		}
		return nil
	case h.Router == nil && len(h.Workers) > 0:
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		h.wg = &sync.WaitGroup{}

		h.work(ctx)

		waitCh := make(chan struct{})
		go func() {
			h.wg.Wait()
			close(waitCh)
		}()

		select {
		case <-ctx.Done():
			logger.Info(context.Background()).Msg("host shutdown started")
			defer logger.Info(context.Background()).Msg("host shutdown completed")
		case <-waitCh:
			logger.Info(context.Background()).Msg("all workers done processing, host is shutting down")
		}

		h.wg.Wait()
		stop()
		return nil
	default:
		return fmt.Errorf("no router or workers added, host will not start")
	}
}

func (h *Host) serve(ctx context.Context, cancel context.CancelFunc) {
	if h.Addr == "" {
		h.Addr = ":8080"
	}

	h.server = &http.Server{
		Addr:    h.Addr,
		Handler: h.Router,
	}

	go func() {
		if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Err(ctx, err).Msg("HTTP server stopped")
			cancel()
		}
	}()
}

func (h *Host) work(ctx context.Context) {
	for _, w := range h.Workers {
		h.wg.Add(1)
		worker := w
		go worker.Run(ctx, h.wg.Done)
	}
}
