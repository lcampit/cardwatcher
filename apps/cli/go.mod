module github.com/lcampit/cardwatcher/apps/cli

go 1.24.0

require (
	github.com/jedib0t/go-pretty/v6 v6.7.8
	github.com/spf13/cobra v1.10.2
	github.com/spf13/viper v1.21.0
	google.golang.org/grpc v1.70.0
	google.golang.org/protobuf v1.36.5
	github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1 v1.0.0
)

replace github.com/lcampit/cardwatcher/gen/go/cardwatcher/v1 => ../../gen/go/cardwatcher/v1
