package service

import (
	"context"
	"final/internal/domain"
	"fmt"
	"go.uber.org/zap"
)

func NewRateService(repo RateSaver, fetcher RateFetcher, logger *zap.SugaredLogger) *RateService {
	return &RateService{
		repo:    repo,
		fetcher: fetcher,
		l:       logger,
	}
}

type RateService struct {
	repo    RateSaver
	fetcher RateFetcher
	l       *zap.SugaredLogger
}

type RateSaver interface {
	SaveRate(ctx context.Context, rate *domain.Rate) error
}

type RateFetcher interface {
	FetchRate(ctx context.Context) (*domain.Rate, error)
}

func (r *RateService) GetRate(ctx context.Context) (*domain.Rate, error) {
	currentRate, err := r.fetcher.FetchRate(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch rate: %w", err)
	}

	err = r.repo.SaveRate(ctx, currentRate)

	if err != nil {
		r.l.Errorf("failed to save rate: %v", err)
	}

	return currentRate, nil
}
