package main

import (
	"fmt"
	"log/slog"
	"net"
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
	"google.golang.org/grpc/reflection"
)

type WatcherConfig struct {
	AppMode              string `env:"APP_MODE"`
	LogLevel             string `env:"LOG_LEVEL"`
	Port                 int    `env:"SERVER_PORT"`
	AccessToken          string `env:"CARDTRADER_ACCESS_TOKEN"`
	MongoHost            string `env:"MONGO_HOST"`
	MongoPort            string `env:"MONGO_PORT"`
	MongoDatabase        string `env:"MONGO_DATABASE"`
	CardtraderAPIBaseURL string `env:"CARDTRADER_API_BASE_URL"`
	NtfyHost             string `env:"NTFY_HOST"`
	NtfyPort             string `env:"NTFY_PORT"`
	NotificationSchedule string `env:"NOTIFICATION_SCHEDULE"`
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

	ntfyAdapter := ntfy.NewNtfyAdapter(logger, "ntfy.sh", "")
	cardtraderAdapter := cardtrader.NewCardtraderAdapter(logger, watcherConfig.AccessToken, watcherConfig.CardtraderAPIBaseURL)
	mongoAdapter := mongo.NewMongoAdapter(logger, watcherConfig.MongoHost, watcherConfig.MongoPort, watcherConfig.MongoDatabase)
	service := service.NewService(logger, cardtraderAdapter, mongoAdapter, ntfyAdapter)
	server := server.NewServer(logger, service)

	grpcServer := grpc.NewServer()
	models.RegisterCardWatcherServer(grpcServer, server)
	reflection.Register(grpcServer)

	loc, _ := time.LoadLocation("Europe/Rome")
	c := cron.New(cron.WithLocation(loc))
	_, err = c.AddFunc(watcherConfig.NotificationSchedule, service.WatchAndNotify)
	if err != nil {
		logger.Error("error when setting up notification cron job", slog.Any("error", err))
	}
	c.Start()
	logger.Info("Server started", slog.Int("serverPort", watcherConfig.Port))
	err = grpcServer.Serve(lis)
	if err != nil {
		logger.Error("error while listening", slog.Any("error", err))
	}
	logger.Info("Stopping server")
	c.Stop()
	logger.Info("Done")
}
