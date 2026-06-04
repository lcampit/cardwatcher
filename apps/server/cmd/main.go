package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lcampit/cardwatcher/apps/server/internal/app"
	"github.com/lcampit/cardwatcher/apps/server/internal/cardtrader"
	"github.com/lcampit/cardwatcher/apps/server/internal/handler"
	"github.com/lcampit/cardwatcher/apps/server/internal/logger"
	"github.com/lcampit/cardwatcher/apps/server/internal/mongo"
	"github.com/lcampit/cardwatcher/apps/server/internal/ntfy"
	"github.com/lcampit/cardwatcher/apps/server/internal/service"

	"go-simpler.org/env"
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
	MongoAuthDatabase                     string `env:"MONGO_AUTH_DATABASE"`
	MongoWatchCollectioName               string `env:"MONGO_WATCH_COLLECTION_NAME"`
	MongoCAFile                           string `env:"MONGO_TLS_CA_FILE"`
	MongoUseReplicaSet                    bool   `env:"MONGO_USE_REPLICA_SET" default:"false"`
	CardtraderAPIBaseURL                  string `env:"CARDTRADER_API_BASE_URL"`
	CardtraderAccessToken                 string `env:"CARDTRADER_ACCESS_TOKEN"`
	NtfyHost                              string `env:"NTFY_HOST"`
	NtfyPort                              string `env:"NTFY_PORT"`
	NtfyTopic                             string `env:"NTFY_TOPIC"`
	HttpSkipVerify                        bool   `env:"HTTP_SKIP_VERIFY"`
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

	cardtraderAdapterConfig := cardtrader.CardtraderAdapterConfig{
		Logger:      logger,
		AccessToken: watcherConfig.CardtraderAccessToken,
		BaseURL:     watcherConfig.CardtraderAPIBaseURL,
		SkipVerify:  watcherConfig.HttpSkipVerify,
	}
	cardtraderAdapter := cardtrader.NewCardtraderAdapter(cardtraderAdapterConfig)

	ntfyAdapterConfig := ntfy.NtfyAdapterConfig{
		Logger: logger,
		Host:   watcherConfig.NtfyHost,
		Port:   watcherConfig.NtfyPort,
		Topic:  watcherConfig.NtfyTopic,
	}
	ntfyAdapter := ntfy.NewNtfyAdapter(ntfyAdapterConfig)

	mongoAdapterConfig := mongo.MongoAdapterConfig{
		Logger:              logger,
		Host:                watcherConfig.MongoHost,
		Port:                watcherConfig.MongoPort,
		Username:            watcherConfig.MongoUsername,
		Password:            watcherConfig.MongoPassword,
		CAFile:              watcherConfig.MongoCAFile,
		Database:            watcherConfig.MongoDatabase,
		AuthDatabase:        watcherConfig.MongoAuthDatabase,
		WatchCollectionName: watcherConfig.MongoWatchCollectioName,
		UseReplicaSet:       watcherConfig.MongoUseReplicaSet,
	}
	mongoAdapter, err := mongo.NewMongoAdapter(mongoAdapterConfig)
	if err != nil {
		logger.Error("error creating mongo adapter", slog.Any("error", err))
		os.Exit(1)
	}

	serviceConfig := service.ServiceConfig{
		Logger:               logger,
		CardtraderAdapter:    cardtraderAdapter,
		MongoAdapter:         mongoAdapter,
		NtfyAdapter:          ntfyAdapter,
		NotificationSchedule: watcherConfig.NotificationSchedule,
		UpdateMapsSchedule:   watcherConfig.UpdateMapsSchedule,
	}
	service := service.NewService(ctx, serviceConfig)
	defer service.Close()

	handlerConfig := handler.HandlerConfig{
		Logger:  logger,
		Service: service,
	}
	handler := handler.NewHandler(handlerConfig)

	app, err := app.NewProductionApp(
		handler,
		logger,
		watcherConfig.ServerPort,
		watcherConfig.ServerEnableReflection,
	)
	if err != nil {
		logger.Error("failed to create cardwatcher server app", slog.Any("error", err))
		os.Exit(1)
	}

	app.StartHealthChecks(time.Duration(watcherConfig.ServerHealthCheckIntervalMilliseconds))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		app.Shutdown(5 * time.Second)
	}()

	if err = app.Run(); err != nil {
		logger.Error("error running app", slog.Any("error", err))
		os.Exit(1)
	}

	os.Exit(0)
}
