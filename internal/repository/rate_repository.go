package repository

import (
	"context"
	"final/internal/domain"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

func NewRateRepository(db *sqlx.DB) *RateRepository {
	return &RateRepository{db: db}
}

type RateRepository struct {
	db *sqlx.DB
}

func (r *RateRepository) SaveRate(ctx context.Context, rate *domain.Rate) error {
	query := `
		INSERT INTO "Rate" ("ask", "bid", "timestamp")
		VALUES ($1, $2, $3)
	`

	_, err := r.db.ExecContext(ctx, query, rate.Ask, rate.Bid, rate.Timestamp.Format(time.RFC3339))

	if err != nil {
		return fmt.Errorf("error while executing SaveRate sql request: %w", err)
	}

	return nil
}
