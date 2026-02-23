module github.com/lcampit/cardwatcher/apps/gateway

go 1.24.0

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.3
	github.com/gorilla/mux v1.8.1
	google.golang.org/grpc v1.70.0
	google.golang.org/protobuf v1.36.5
	github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1 v1.0.0
)

replace github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1 => ../../gen/go/cardwatcher/v1
