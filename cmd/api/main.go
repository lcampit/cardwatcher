package main

import (
	"card-watcher/internal/cardtrader"
	"card-watcher/internal/models"
	"card-watcher/internal/mongo"
	"card-watcher/internal/server"
	"card-watcher/internal/service"
	"fmt"
	"net"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go-simpler.org/env"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type WatcherConfig struct {
	AppMode              string `env:"APP_MODE"`
	Port                 int    `env:"PORT"`
	AccessToken          string `env:"CARDTRADER_ACCESS_TOKEN"`
	MongoHost            string `env:"MONGO_HOST"`
	MongoPort            string `env:"MONGO_PORT"`
	MongoDatabase        string `env:"MONGO_DATABASE"`
	CardtraderApiBaseUrl string `env:"CARDTRADER_API_BASE_URL"`
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

	cardtraderAdapter := cardtrader.NewCardtraderAdapter(watcherConfig.AccessToken, watcherConfig.CardtraderApiBaseUrl)
	mongoAdapter := mongo.NewMongoAdapter(watcherConfig.MongoHost, watcherConfig.MongoPort, watcherConfig.MongoDatabase)
	service := service.NewService(cardtraderAdapter, mongoAdapter)
	server := server.NewServer(service)

	grpcServer := grpc.NewServer()
	models.RegisterCardWatcherServer(grpcServer, server)
	reflection.Register(grpcServer)

	log.Info().Msgf("Server started on port %d", watcherConfig.Port)
	grpcServer.Serve(lis)
}
