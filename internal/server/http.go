package server

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	v1 "sensor/api/sensor/v1"
	"sensor/internal/conf"
	"sensor/internal/service"

	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer news an HTTP sensor. For gRPC requests, refer to the function [NewGRPCServer].
//
// This function would read the configuration to configure the HTTP sensor well,
// and then register the service to the HTTP sensor.
func NewHTTPServer(c *conf.Server, s *service.SensorService, m Middlewares) *http.Server {
	// Here we tell the framework that we need these middlewares, and the framework would provide them automatically.
	opts := []http.ServerOption{
		http.Middleware(m...),
	}
	// Reflect the options on the configuration into the memory
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	// Instantiate a new HTTP sensor listening to a specific port to serve the requests.
	srv := http.NewServer(opts...)
	srv.Handle("/metrics", promhttp.Handler())  // We shall register the Prometheus handler to the sensor as well
	v1.RegisterUserManagementHTTPServer(srv, s) // Register the service handlers as well
	return srv
}
