package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lcampit/cardwatcher/apps/gateway/internal/logger"
	apiv1 "github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go-simpler.org/env"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GatewayConfig struct {
	LogLevel    string `env:"LOG_LEVEL" default:"debug"`
	GatewayPort int    `env:"GATEWAY_PORT"`
	ServerHost  string `env:"SERVER_HOST"`
	ServerPort  int    `env:"SERVER_PORT"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var gatewayConfig GatewayConfig
	if err := env.Load(&gatewayConfig, nil); err != nil {
		fmt.Println(err)
		return
	}

	logger := logger.NewLogger(gatewayConfig.LogLevel)

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	serverAddr := fmt.Sprintf("%s:%d", gatewayConfig.ServerHost, gatewayConfig.ServerPort)
	err := apiv1.RegisterCardWatcherServiceHandlerFromEndpoint(ctx, mux, serverAddr, opts)
	if err != nil {
		logger.Error("error registering gateway", slog.Any("error", err))
		os.Exit(1)
	}

	go func() {
		logger.Info("starting REST gateway",
			slog.Int("port", gatewayConfig.GatewayPort),
			slog.String("server", serverAddr))
		server := http.Server{
			Addr:              serverAddr,
			Handler:           mux,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
		}
		err = server.ListenAndServe()
		if err != nil {
			logger.Error("error starting gateway", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	logger.Info("gateway started", slog.Int("port", gatewayConfig.GatewayPort))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	logger.Info("stopping gateway")
	logger.Info("done")
	os.Exit(0)
}
