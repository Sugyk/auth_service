package grpc_api

import (
	"context"
	"fmt"
	"net"

	"github.com/Sugyk/auth_service/internal/api/grpc/pb"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type Router struct {
	server   *grpc.Server
	listener net.Listener
}

func NewRouter(addr string, authServer pb.AuthServiceServer) (*Router, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("listen on %s: %w", addr, err)
	}

	srvMetrics := grpcprom.NewServerMetrics(grpcprom.WithServerHandlingTimeHistogram())
	prometheus.MustRegister(srvMetrics)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(srvMetrics.UnaryServerInterceptor()),
		grpc.ChainStreamInterceptor(srvMetrics.StreamServerInterceptor()),
	)

	pb.RegisterAuthServiceServer(server, authServer)

	healthServer := health.NewServer()
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(server, healthServer)

	reflection.Register(server)

	srvMetrics.InitializeMetrics(server)

	return &Router{
		server:   server,
		listener: listener,
	}, nil
}

func (r *Router) Start() error {
	return r.server.Serve(r.listener)
}

func (r *Router) Shutdown(ctx context.Context) error {
	stopped := make(chan struct{})
	go func() {
		r.server.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		return nil
	case <-ctx.Done():
		r.server.Stop()
		return ctx.Err()
	}
}
