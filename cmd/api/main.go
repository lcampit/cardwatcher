package main

import (
	"card-watcher/internal/cardtrader"
	"card-watcher/internal/mongo"
	"card-watcher/internal/watcher"
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"go-simpler.org/env"
)

type WatcherConfig struct {
	Port                 int    `env:"PORT"`
	AccessToken          string `env:"ACCESS_TOKEN"`
	MongoHost            string `env:"MONGO_HOST"`
	MongoPort            string `env:"MONGO_PORT"`
	MongoDatabase        string `env:"MONGO_DATABASE"`
	CardtraderApiBaseUrl string `env:"CARDTRADER_API_BASE_URL"`
}

func gracefulShutdown(fiberServer *watcher.Watcher, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := fiberServer.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {
	var watcherConfig WatcherConfig
	if err := env.Load(&watcherConfig, nil); err != nil {
		fmt.Println(err)
		return
	}

	cardtraderAdapter := cardtrader.NewCardtraderAdapter(watcherConfig.AccessToken, watcherConfig.CardtraderApiBaseUrl)
	mongoAdapter := mongo.NewMongoAdapter(watcherConfig.MongoHost, watcherConfig.MongoPort, watcherConfig.MongoDatabase)

	watcher := watcher.New(cardtraderAdapter, mongoAdapter)
	watcher.RegisterFiberRoutes()

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	go func() {
		err := watcher.Listen(fmt.Sprintf(":%d", watcherConfig.Port))
		if err != nil {
			panic(fmt.Sprintf("http server error: %s", err))
		}
	}()

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(watcher, done)

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}
