package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"card-watcher/internal/cardtrader"
	"card-watcher/internal/logger"
	"card-watcher/internal/models"
	"card-watcher/internal/mongo"
	"card-watcher/internal/ntfy"
	"card-watcher/internal/server"
	"card-watcher/internal/service"

	"github.com/robfig/cron/v3"
	"go-simpler.org/env"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type WatcherConfig struct {
	LogLevel                        string `env:"LOG_LEVEL"`
	Port                            int    `env:"SERVER_PORT"`
	AccessToken                     string `env:"CARDTRADER_ACCESS_TOKEN"`
	MongoHost                       string `env:"MONGO_HOST"`
	MongoPort                       string `env:"MONGO_PORT"`
	MongoDatabase                   string `env:"MONGO_DATABASE"`
	MongoWatchCollectioName         string `env:"MONGO_WATCH_COLLECTION_NAME"`
	CardtraderAPIBaseURL            string `env:"CARDTRADER_API_BASE_URL"`
	NtfyHost                        string `env:"NTFY_HOST"`
	NtfyPort                        string `env:"NTFY_PORT"`
	NotificationSchedule            string `env:"NOTIFICATION_SCHEDULE"`
	HeatlhCheckIntervalMilliseconds int    `env:"HEALTH_CHECK_INTERVAL_MILLISECONDS"`
}

func main() {
	var watcherConfig WatcherConfig
	if err := env.Load(&watcherConfig, nil); err != nil {
		fmt.Println(err)
		return
	}

	logger := logger.NewLogger(watcherConfig.LogLevel)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", watcherConfig.Port))
	if err != nil {
		logger.Error("failed to listen", slog.Any("error", err))
	}

	logger.Info("creating cardtrader adapter")
	cardtraderAdapterConfig := cardtrader.CardtraderAdapterConfig{
		Logger:      logger,
		AccessToken: watcherConfig.AccessToken,
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
	}
	mongoAdapter, err := mongo.NewMongoAdapter(mongoAdapterConfig)
	if err != nil {
		logger.Error("error creating mongo adapter", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("creating service")
	serviceConfig := service.ServiceConfig{
		Logger:            logger,
		CardtraderAdapter: cardtraderAdapter,
		MongoAdapter:      mongoAdapter,
		NtfyAdapter:       ntfyAdapter,
	}
	service := service.NewService(serviceConfig)

	logger.Info("creating server")
	serverConfig := server.ServerConfig{
		Logger:  logger,
		Service: service,
	}
	server := server.NewServer(serverConfig)

	grpcServer := grpc.NewServer()
	healthcheck := health.NewServer()
	healthgrpc.RegisterHealthServer(grpcServer, healthcheck)
	models.RegisterCardWatcherServer(grpcServer, server)
	reflection.Register(grpcServer)

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
			time.Sleep(time.Duration(watcherConfig.HeatlhCheckIntervalMilliseconds) * time.Millisecond)
		}
	}()

	loc, _ := time.LoadLocation("Europe/Rome")
	c := cron.New(cron.WithLocation(loc))
	_, err = c.AddFunc(watcherConfig.NotificationSchedule, service.WatchAndNotify)
	if err != nil {
		logger.Error("error when setting up notification cron job", slog.Any("error", err))
	}
	c.Start()
	logger.Info("server started", slog.Int("serverPort", watcherConfig.Port))
	err = grpcServer.Serve(lis)
	if err != nil {
		logger.Error("error while listening", slog.Any("error", err))
	}
	logger.Info("stopping server")
	c.Stop()
	logger.Info("done")
}
