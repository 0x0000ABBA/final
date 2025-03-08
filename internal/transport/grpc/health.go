package grpc

import (
	"context"
	"final/internal/transport/gen"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func NewHealthServiceServer() *HealthServiceServer {
	tracer := otel.Tracer("final-service/health")
	return &HealthServiceServer{
		tracer: tracer,
	}
}

type HealthServiceServer struct {
	gen.UnimplementedHealthServiceServer
	tracer trace.Tracer
}

func (s *HealthServiceServer) HealthCheck(ctx context.Context, req *gen.HealthCheckRequest) (*gen.HealthCheckResponse, error) {
	_, span := s.tracer.Start(ctx, "HealthCheck")

	res := &gen.HealthCheckResponse{OK: true}

	span.End()

	return res, nil
}
