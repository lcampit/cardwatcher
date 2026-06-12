// Package app handles the grpc server lifecycle,
// leaving main free to just initialize the app
// and running it
//
// App handles the server start, its graceful stop
// and healthcheck logic
package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"buf.build/go/protovalidate"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/lcampit/cardwatcher/apps/server/internal/handler"
	"github.com/lcampit/cardwatcher/apps/server/internal/handler/middleware"
	apiv1 "github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1"
)

type App struct {
	grpcServer   *grpc.Server
	healthServer *health.Server
	logger       *slog.Logger
	listener     net.Listener
	handler      handler.Handler
	ctx          context.Context
	cancel       context.CancelFunc
}

// NewProductionApp initializes and registers the grpc and health servers
// with a production ready tcp listener
func NewProductionApp(
	handler handler.Handler,
	logger *slog.Logger,
	serverPort int,
	enableReflection bool,
) (*App, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", serverPort))
	if err != nil {
		return nil, err
	}
	return NewApp(handler, logger, lis, serverPort, enableReflection)
}

// NewApp initializes and registers the grpc and health servers
// with the listener provided
func NewApp(
	handler handler.Handler,
	logger *slog.Logger,
	listener net.Listener,
	serverPort int,
	enableReflection bool,
) (*App, error) {
	ctx, cancel := context.WithCancel(context.Background())

	validator, err := protovalidate.New()
	if err != nil {
		cancel()
		return nil, err
	}

	// all interceptors are chained here
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			middleware.LoggingInterceptor(logger),
			protovalidate_middleware.UnaryServerInterceptor(validator),
			middleware.ErrorInterceptor,
		),
	}

	grpcServer := grpc.NewServer(opts...)
	apiv1.RegisterCardWatcherServiceServer(grpcServer, handler)

	if enableReflection {
		reflection.Register(grpcServer)
	}

	healthcheck := health.NewServer()
	healthgrpc.RegisterHealthServer(grpcServer, healthcheck)
	healthcheck.SetServingStatus("", healthgrpc.HealthCheckResponse_SERVING)

	return &App{
		grpcServer:   grpcServer,
		healthServer: healthcheck,
		logger:       logger,
		handler:      handler,
		listener:     listener,
		ctx:          ctx,
		cancel:       cancel,
	}, nil
}

// StartHealthChecks runs the background loop to update health status.
func (a *App) StartHealthChecks(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-a.ctx.Done():
				return
			case <-ticker.C:
				a.updateHealthStatus()
			}
		}
	}()
}

func (a *App) updateHealthStatus() {
	ctx, cancel := context.WithTimeout(a.ctx, 2*time.Second)
	defer cancel()

	if err := a.handler.Health(ctx); err != nil {
		a.logger.Error("error in healthcheck", slog.Any("error", err))
		a.healthServer.SetServingStatus("", healthgrpc.HealthCheckResponse_NOT_SERVING)
		return
	}
	a.healthServer.SetServingStatus("", healthgrpc.HealthCheckResponse_SERVING)
}

// Run starts the server and blocks.
func (a *App) Run() error {
	a.logger.Info("starting grpc server", slog.String("addr", a.listener.Addr().String()))
	return a.grpcServer.Serve(a.listener)
}

// Shutdown handles the graceful termination sequence.
func (a *App) Shutdown(timeout time.Duration) {
	a.logger.Info("initiating graceful shutdown")

	// Mark NOT_SERVING immediately to fail readiness probes
	a.healthServer.SetServingStatus("", healthgrpc.HealthCheckResponse_NOT_SERVING)

	// Wait for load balancers to drain (critical for K8s)
	time.Sleep(5 * time.Second)

	// Graceful stop with timeout
	done := make(chan struct{})
	go func() {
		a.grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		a.logger.Debug("graceful shutdown completed")
	case <-time.After(timeout):
		a.logger.Debug("forced shutdown due to timeout")
		a.grpcServer.Stop()
	}

	a.cancel()
}
