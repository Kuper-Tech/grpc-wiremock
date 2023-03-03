package health

import (
	"context"

	"google.golang.org/grpc"

	"grpc-proxy/pkg/gitlab.sbmt.io/paas/health"
)

type Service struct {
	health.UnimplementedHealthServer
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) RegisterGRPC(server *grpc.Server) {
	health.RegisterHealthServer(server, s)
}

func (s *Service) Check(_ context.Context, req *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	return &health.HealthCheckResponse{
		Status: health.HealthCheckResponse_SERVING,
	}, nil
}

func (s *Service) Watch(_ *health.HealthCheckRequest, server health.Health_WatchServer) error {
	return server.Send(&health.HealthCheckResponse{
		Status: health.HealthCheckResponse_SERVING,
	})
}
