package main

import (
	"flag"
	"fmt"
	"github.com/go-kratos/kratos/v2/registry"
	"os"
	"strings"
	"time"

	"sensor/internal/conf"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	_ "go.uber.org/automaxprocs"
)

// The process of initialization of the service follows the route below:
//
// `var` defined global variables -> init() function -> main() function
//
// If you have never coded with Golang, you should read the comments in this file carefully.
// We recommend you comment your work in English so that collaborators could grasp your idea directly
// from your friendly comments.
//
// Copyright (c) 2024 Ryker Zhu

// You can specify the name and version of your service by passing a special parameter '-ldflags'
// during Golang compilation, just as what the following command shows:
//
//	go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	// Each service should have its own unique name
	Name = "example-service"
	// Version is the version of the compiled software.
	// You should name the version as the following format:
	// major version . minor version . patch number
	Version = "1.0.0"

	// flagconf is the config flag, which specifies the location of the configuration file
	flagconf string
	// We adopt the hostname as the identifier of the service
	id, _ = os.Hostname()

	startTime  = time.Now()
	timeFormat = "2006-01-02 15:04:05.000"
)

// init initializes the service before the main function executes, so you should place all the initialization
// code inside this function. For instance, you can declare more flags or introduce more complicated mechanism
// into the code.
func init() {
	// Declare a flag '-conf' so that the user can pass a config file with the argument
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

// newApp initializes all the dependencies that the service requires and then creates the service instance
//
// Required dependency instances are passed to the function automatically by the dependency injection, so
// you should place all the dependencies required in the parentheses as parameters and during compilation,
// the code of injecting the dependencies would be generated by wire.
//
// We decided to use the following combo of technologies:
//
//   - Service discovery & registration: [etcd]: https://github.com/etcd-io/etcd
//   - Logging: [zap]: https://github.com/uber-go/zap
//   - HTTP router: [gin]: https://github.com/gin-gonic/gin
//   - Path tracing: [OpenTelemetry]: http://opentelemetry.io/
//   - Object relation mapping: [Ent]: https://entgo.io/
//   - Rate limiter: Kratos builtin middleware
//   - Circuit breaker: Kratos builtin middleware
//   - Authorization: Json Web Token
//
// DO NOT HARD CODE CONFIG OR DEPENDENCIES
func newApp(reg registry.Registrar, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),           // A service ID should be unique in the global scope
		kratos.Name(Name),       // A service name should be human-readable and clear enough to ensure maintainability
		kratos.Version(Version), // Current version of this service
		kratos.Metadata(map[string]string{ // Provides metadata about this service
			"description": "The service provides a basic microservice framework.",
		}),
		kratos.Server( // The service runs both HTTP and GRPC sensor simultaneously.
			gs, hs, // Intro-service calls should utilize GRPC sensor while the front end uses HTTP sensor
		),
		kratos.Registrar(reg), // Tell the Kratos to use the client as its registrar
	)
}

// main function is in charge of loading the configuration and then starts the service
func main() {
	// Parse arguments from the command line
	flag.Parse()
	// Initialize the config with the source
	// Note that a source may on the local disk, but it can also be fetched from remote one
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer func(c config.Config) {
		err := c.Close()
		if err != nil {
			log.Warn(err)
		}
	}(c) // We should make sure the config file would finally be closed

	if err := c.Load(); err != nil { // Try to load the config file from the given source
		panic(err)
	}

	// Bootstrap is the configuration structure of the config file
	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil { // Try to deserialize the config into the corresponding structure
		panic(err)
	}

	log.SetLogger(NewLogger(bc.Telemetry.Log))

	// Inject dependencies into the service
	app, cleanup, err := wireApp(bc.Registry, bc.Server, bc.Data, bc.Telemetry)
	if err != nil {
		panic(err)
	}
	defer cleanup() // Clean up the injected dependencies before exits

	builder := strings.Builder{}
	builder.WriteString(
		fmt.Sprintf(
			"\n    \033[36;1m%s\033[0m \033[96m%s\033[0m  ready in \033[33;1m%v\033[0m\n",
			Name, Version, time.Now().Sub(startTime),
		),
	)
	fmt.Println(builder.String())

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}