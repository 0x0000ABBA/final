package grpc

import (
	"context"
	"final/internal/domain"
	"final/internal/transport/gen"
	"fmt"
)

func NewRateServiceServer(service RateService) *RateServiceServer {
	if service == nil {
		return nil
	}
	return &RateServiceServer{
		service: service,
	}
}

type RateServiceServer struct {
	service RateService
	gen.UnimplementedRateServiceServer
}

type RateService interface {
	GetRate(ctx context.Context) (*domain.Rate, error)
}

func (s *RateServiceServer) GetRate(ctx context.Context, req *gen.GetRateRequest) (*gen.GetRateResponse, error) {
	rate, err := s.service.GetRate(ctx)

	if err != nil {
		return nil, fmt.Errorf("error while using rate service: %w", err)
	}

	return &gen.GetRateResponse{
		Ask:       rate.Ask,
		Bid:       rate.Bid,
		Timestamp: rate.Timestamp.String(),
	}, nil
}
