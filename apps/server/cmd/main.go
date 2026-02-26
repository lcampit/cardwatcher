package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lcampit/cardwatcher/apps/server/internal/cardtrader"
	"github.com/lcampit/cardwatcher/apps/server/internal/handler"
	"github.com/lcampit/cardwatcher/apps/server/internal/logger"
	"github.com/lcampit/cardwatcher/apps/server/internal/mongo"
	"github.com/lcampit/cardwatcher/apps/server/internal/ntfy"
	"github.com/lcampit/cardwatcher/apps/server/internal/service"
	apiv1 "github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1"

	"go-simpler.org/env"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type WatcherConfig struct {
	LogLevel                              string `env:"LOG_LEVEL" default:"debug"`
	ServerPort                            int    `env:"SERVER_PORT"`
	NotificationSchedule                  string `env:"NOTIFICATION_SCHEDULE"`
	UpdateMapsSchedule                    string `env:"UPDATE_MAPS_SCHEDULE"`
	ServerHealthCheckIntervalMilliseconds int    `env:"SERVER_HEALTH_CHECK_INTERVAL_MILLISECONDS" default:"1000"`
	ServerEnableReflection                bool   `env:"SERVER_ENABLE_REFLECTION" default:"false"`
	MongoHost                             string `env:"MONGO_HOST"`
	MongoPort                             string `env:"MONGO_PORT"`
	MongoUsername                         string `env:"MONGO_USERNAME"`
	MongoPassword                         string `env:"MONGO_PASSWORD"`
	MongoDatabase                         string `env:"MONGO_DATABASE"`
	MongoWatchCollectioName               string `env:"MONGO_WATCH_COLLECTION_NAME"`
	MongoCAFile                           string `env:"MONGO_TLS_CA_FILE"`
	MongoUseReplicaSet                    bool   `env:"MONGO_USE_REPLICA_SET" default:"false"`
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
		Username:            watcherConfig.MongoUsername,
		Password:            watcherConfig.MongoPassword,
		CAFile:              watcherConfig.MongoCAFile,
		Database:            watcherConfig.MongoDatabase,
		WatchCollectionName: watcherConfig.MongoWatchCollectioName,
		UseReplicaSet:       watcherConfig.MongoUseReplicaSet,
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
	handlerConfig := handler.HandlerConfig{
		Logger:  logger,
		Service: service,
	}
	handler := handler.NewHandler(handlerConfig)

	grpcServer := grpc.NewServer()
	healthcheck := health.NewServer()
	healthgrpc.RegisterHealthServer(grpcServer, healthcheck)
	apiv1.RegisterCardWatcherServiceServer(grpcServer, handler)

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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		err = grpcServer.Serve(lis)
		if err != nil {
			logger.Error("error while listening", slog.Any("error", err))
		}
	}()
	logger.Info("server started", slog.Int("serverPort", watcherConfig.ServerPort))
	// Waiting for cancel signal
	<-c
	logger.Info("stopping server")
	// Gracefulstop will block until current requests are completed
	grpcServer.GracefulStop()
	service.Close()

	logger.Info("done")
	os.Exit(0)
}
