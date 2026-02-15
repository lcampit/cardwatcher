package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	api "github.com/lcampit/card-watcher-server/internal/api/v1"
	"github.com/lcampit/card-watcher-server/internal/server/cardtrader"
	"github.com/lcampit/card-watcher-server/internal/server/logger"
	"github.com/lcampit/card-watcher-server/internal/server/mongo"
	"github.com/lcampit/card-watcher-server/internal/server/ntfy"
	"github.com/lcampit/card-watcher-server/internal/server/server"
	"github.com/lcampit/card-watcher-server/internal/server/service"

	"go-simpler.org/env"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type WatcherConfig struct {
	LogLevel                              string `env:"LOG_LEVEL"`
	ServerPort                            int    `env:"SERVER_PORT"`
	NotificationSchedule                  string `env:"NOTIFICATION_SCHEDULE"`
	UpdateMapsSchedule                    string `env:"UPDATE_MAPS_SCHEDULE"`
	ServerHealthCheckIntervalMilliseconds int    `env:"SERVER_HEALTH_CHECK_INTERVAL_MILLISECONDS" default:"1000"`
	ServerEnableReflection                bool   `env:"SERVER_ENABLE_REFLECTION" default:"false"`
	MongoHost                             string `env:"MONGO_HOST"`
	MongoPort                             string `env:"MONGO_PORT"`
	MongoDatabase                         string `env:"MONGO_DATABASE"`
	MongoWatchCollectioName               string `env:"MONGO_WATCH_COLLECTION_NAME"`
	MongoConnectionRetries                int    `env:"MONGO_CONNECTION_RETRIED" default:"5"`
	CardtraderAPIBaseURL                  string `env:"CARDTRADER_API_BASE_URL"`
	CardtraderAccessToken                 string `env:"CARDTRADER_ACCESS_TOKEN"`
	NtfyHost                              string `env:"NTFY_HOST"`
	NtfyPort                              string `env:"NTFY_PORT"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var watcherConfig WatcherConfig
	if err := env.Load(&watcherConfig, nil); err != nil {
		fmt.Println(err)
		return
	}

	logger := logger.NewLogger(watcherConfig.LogLevel)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", watcherConfig.ServerPort))
	if err != nil {
		logger.Error("failed to listen", slog.Any("error", err))
	}

	logger.Info("creating cardtrader adapter")
	cardtraderAdapterConfig := cardtrader.CardtraderAdapterConfig{
		Logger:      logger,
		AccessToken: watcherConfig.CardtraderAccessToken,
		BaseURL:     watcherConfig.CardtraderAPIBaseURL,
	}
	cardtraderAdapter := cardtrader.NewCardtraderAdapter(cardtraderAdapterConfig)

	logger.Info("creating ntfy adapter")
	ntfyAdapterConfig := ntfy.NtfyAdapterConfig{
		Logger: logger,
		Host:   watcherConfig.NtfyHost,
		Port:   watcherConfig.NtfyPort,
	}
	ntfyAdapter := ntfy.NewNtfyAdapter(ntfyAdapterConfig)

	logger.Info("creating mongo adapter")
	mongoAdapterConfig := mongo.MongoAdapterConfig{
		Logger:              logger,
		Host:                watcherConfig.MongoHost,
		Port:                watcherConfig.MongoPort,
		Database:            watcherConfig.MongoDatabase,
		WatchCollectionName: watcherConfig.MongoWatchCollectioName,
		ConnectionRetries:   watcherConfig.MongoConnectionRetries,
	}
	mongoAdapter, err := mongo.NewMongoAdapter(mongoAdapterConfig)
	if err != nil {
		logger.Error("error creating mongo adapter", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("creating service")
	serviceConfig := service.ServiceConfig{
		Logger:               logger,
		CardtraderAdapter:    cardtraderAdapter,
		MongoAdapter:         mongoAdapter,
		NtfyAdapter:          ntfyAdapter,
		NotificationSchedule: watcherConfig.NotificationSchedule,
		UpdateMapsSchedule:   watcherConfig.UpdateMapsSchedule,
	}
	service := service.NewService(ctx, serviceConfig)

	logger.Info("creating server")
	serverConfig := server.ServerConfig{
		Logger:  logger,
		Service: service,
	}
	server := server.NewServer(serverConfig)

	grpcServer := grpc.NewServer()
	healthcheck := health.NewServer()
	healthgrpc.RegisterHealthServer(grpcServer, healthcheck)
	api.RegisterCardWatcherServer(grpcServer, server)

	if watcherConfig.ServerEnableReflection {
		reflection.Register(grpcServer)
	}

	// Periodically check adapters health
	go func() {
		for {
			err := mongoAdapter.Health()
			if err != nil {
				logger.Error("error in mongo adapter health check", slog.Any("error", err))
				healthcheck.SetServingStatus("", healthgrpc.HealthCheckResponse_NOT_SERVING)
			} else {
				healthcheck.SetServingStatus("", healthgrpc.HealthCheckResponse_SERVING)
			}
			time.Sleep(time.Duration(watcherConfig.ServerHealthCheckIntervalMilliseconds) * time.Millisecond)
		}
	}()

	logger.Info("server started", slog.Int("serverPort", watcherConfig.ServerPort))
	err = grpcServer.Serve(lis)
	if err != nil {
		logger.Error("error while listening", slog.Any("error", err))
	}
	logger.Info("stopping server")
	service.Close()
	logger.Info("done")
}
