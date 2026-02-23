package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GatewayConfig struct {
	ServerPort int    `env:"GATEWAY_PORT"`
	GRPCServer string `env:"GRPC_SERVER"`
	GRPCPort   int    `env:"GRPC_PORT"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := GatewayConfig{
		ServerPort: 8080,
		GRPCServer: "localhost",
		GRPCPort:   50051,
	}

	// Create gRPC connection to server
	grpcServerAddress := fmt.Sprintf("%s:%d", config.GRPCServer, config.GRPCPort)
	conn, err := grpc.NewClient(grpcServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial server: %v", err)
	}
	defer conn.Close()

	// Create gRPC-Gateway mux
	gwmux := runtime.NewServeMux()

	// Register the gateway handlers
	if err := v1.RegisterCardWatcherHandler(ctx, gwmux, conn); err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	// Create HTTP server
	gatewayAddr := fmt.Sprintf(":%d", config.ServerPort)
	gwServer := &http.Server{
		Addr:    gatewayAddr,
		Handler: gwmux,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Starting gateway server on %s", gatewayAddr)
		if err := gwServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gateway server...")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := gwServer.Shutdown(ctx); err != nil {
		log.Fatalf("gateway forced to shutdown: %v", err)
	}

	log.Println("Gateway server exited")
}
