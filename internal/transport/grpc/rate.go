package grpc

import (
	"context"
	"final/internal/domain"
	"final/internal/transport/gen"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	rateRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rate_service_requests_total",
			Help: "Total number of RateService requests",
		},
		[]string{"method"},
	)
	rateErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rate_service_errors_total",
			Help: "Total number of RateService errors",
		},
		[]string{"method"},
	)
)

func init() {
	prometheus.MustRegister(rateRequests, rateErrors)
}

func NewRateServiceServer(service RateService) *RateServiceServer {
	if service == nil {
		return nil
	}

	tracer := otel.Tracer("final-service/rate")

	return &RateServiceServer{
		tracer:  tracer,
		service: service,
	}
}

type RateServiceServer struct {
	tracer  trace.Tracer
	service RateService
	gen.UnimplementedRateServiceServer
}

type RateService interface {
	GetRate(ctx context.Context) (*domain.Rate, error)
}

func (s *RateServiceServer) GetRate(ctx context.Context, req *gen.GetRateRequest) (*gen.GetRateResponse, error) {
	rateRequests.WithLabelValues("GetRate").Inc()
	ctx, span := s.tracer.Start(ctx, "GetRate")

	rate, err := s.service.GetRate(ctx)

	if err != nil {
		rateErrors.WithLabelValues("GetRate").Inc()
		return nil, fmt.Errorf("error while using rate service: %w", err)
	}

	span.End()

	return &gen.GetRateResponse{
		Ask:       rate.Ask,
		Bid:       rate.Bid,
		Timestamp: rate.Timestamp.String(),
	}, nil
}
