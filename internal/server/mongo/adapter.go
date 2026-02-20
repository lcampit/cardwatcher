// Package mongo exposes database related operations
// using a single adapter
package mongo

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoAdapter interface {
	SaveWatch(ctx context.Context, watch *Watch) (string, error)
	GetWatches(ctx context.Context) ([]*Watch, error)
	GetWatchByWatchID(ctx context.Context, watchID string) (*Watch, error)
	DeleteWatchByID(ctx context.Context, watchID string) error

	Health() error
}

type mongoAdapter struct {
	logger          *slog.Logger
	client          *mongo.Client
	database        string
	watchCollection string
}

type MongoAdapterConfig struct {
	Logger              *slog.Logger
	Host                string
	Port                string
	Username            string
	Password            string
	Database            string
	WatchCollectionName string
	CAFile              string
}

func NewMongoAdapter(config MongoAdapterConfig) (MongoAdapter, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	credentials := options.Credential{
		AuthSource:    config.Database,
		Username:      config.Username,
		Password:      config.Password,
		AuthMechanism: "SCRAM-SHA-256",
	}

	host := config.Host
	if host == "" {
		host = "localhost"
	}

	port := config.Port
	if port == "" {
		port = "27017"
	}

	clientOpts := options.Client().
		SetHosts([]string{fmt.Sprintf("%s:%s", host, port)}).
		SetAuth(credentials).
		SetRetryReads(true).
		SetRetryWrites(true).
		SetServerSelectionTimeout(5 * time.Second).
		SetConnectTimeout(10 * time.Second).
		SetMaxPoolSize(50).
		SetMinPoolSize(5)

	if isLocalHost(config.Host) {
		clientOpts.SetDirect(true)
	} else {
		clientOpts.SetReplicaSet("rs0")
	}

	if strings.TrimSpace(config.CAFile) != "" {
		tlsCfg, err := buildTLSConfig(config)
		if err != nil {
			return nil, fmt.Errorf("build TLS config: %w", err)
		}
		clientOpts.SetTLSConfig(tlsCfg)
	}

	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, fmt.Errorf("creating mongo client: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("pinging mongo instance: %w", err)
	}

	return &mongoAdapter{
		logger:          config.Logger,
		client:          client,
		database:        config.Database,
		watchCollection: config.WatchCollectionName,
	}, nil
}

func isLocalHost(h string) bool {
	hl := strings.ToLower(strings.TrimSpace(h))
	return hl == "localhost" || hl == "127.0.0.1"
}

func buildTLSConfig(config MongoAdapterConfig) (*tls.Config, error) {
	caPEM, err := os.ReadFile(config.CAFile)
	if err != nil {
		return nil, fmt.Errorf("reading CA file: %w", err)
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caPEM) {
		return nil, errors.New("appending CA PEM failed")
	}

	tlsConfig := &tls.Config{
		RootCAs:    pool,
		MinVersion: tls.VersionTLS12,
		ServerName: config.Host,
	}

	return tlsConfig, nil
}

func (a *mongoAdapter) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return a.client.Ping(ctx, nil)
}
