module github.com/lcampit/cardwatcher/apps/server

go 1.24.0

require (
	github.com/carlmjohnson/requests v0.24.3
	github.com/robfig/cron/v3 v3.0.1
	go-simpler.org/env v0.12.0
	go.mongodb.org/mongo-driver/v2 v2.2.2
	google.golang.org/grpc v1.70.0
	google.golang.org/grpc/reflection v1.70.0
	google.golang.org/protobuf v1.36.5
	github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1 v1.0.0
)

replace github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1 => ../../gen/go/cardwatcher/v1
