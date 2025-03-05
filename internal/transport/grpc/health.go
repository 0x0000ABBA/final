package grpc

import (
	"context"
	"final/internal/transport/gen"
)

func NewHealthServiceServer() *HealthServiceServer {
	return &HealthServiceServer{}
}

type HealthServiceServer struct {
	gen.UnimplementedHealthServiceServer
}

func (s *HealthServiceServer) HealthCheck(ctx context.Context, req *gen.HealthCheckRequest) (*gen.HealthCheckResponse, error) {
	return &gen.HealthCheckResponse{OK: true}, nil
}
