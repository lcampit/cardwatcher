package main

import (
	"fmt"
	"net"
	"time"

	"card-watcher/internal/cardtrader"
	"card-watcher/internal/models"
	"card-watcher/internal/mongo"
	"card-watcher/internal/ntfy"
	"card-watcher/internal/server"
	"card-watcher/internal/service"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go-simpler.org/env"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type WatcherConfig struct {
	AppMode              string `env:"APP_MODE"`
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

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	if watcherConfig.AppMode == "PROD" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", watcherConfig.Port))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}

	ntfyAdapter := ntfy.NewNtfyAdapter("ntfy.sh", "")
	cardtraderAdapter := cardtrader.NewCardtraderAdapter(watcherConfig.AccessToken, watcherConfig.CardtraderAPIBaseURL)
	mongoAdapter := mongo.NewMongoAdapter(watcherConfig.MongoHost, watcherConfig.MongoPort, watcherConfig.MongoDatabase)
	service := service.NewService(cardtraderAdapter, mongoAdapter, ntfyAdapter)
	server := server.NewServer(service)

	grpcServer := grpc.NewServer()
	models.RegisterCardWatcherServer(grpcServer, server)
	reflection.Register(grpcServer)

	loc, _ := time.LoadLocation("Europe/Rome")
	c := cron.New(cron.WithLocation(loc))
	_, err = c.AddFunc(watcherConfig.NotificationSchedule, service.WatchAndNotify)
	if err != nil {
		log.Fatal().Err(err).Msg("error when setting up notification cron job")
	}
	c.Start()
	log.Info().Msgf("Server started on port %d", watcherConfig.Port)
	grpcServer.Serve(lis)
	log.Info().Msgf("Stopping server")
	c.Stop()
	log.Info().Msgf("Done")
}
