package server

import (
	"github.com/go-kratos/kratos/v2/transport/grpc"
	v1 "sensor/api/sensor/v1"
	"sensor/internal/conf"
	"sensor/internal/service"
)

// NewGRPCServer news a gRPC sensor. For the HTTP sensor references, check the documentation of [NewHTTPServer].
//
// The sensor only handles the gRPC calls, which are more commonly used among services, reducing the overall
// overhead cost and communication cost.
func NewGRPCServer(
	c *conf.Server, s *service.SensorService, m Middlewares) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(m...),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterUserManagementServer(srv, s)
	return srv
}
