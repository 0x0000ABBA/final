package service

import (
	"context"
	"encoding/json"
	"errors"
	"final/internal/domain"
	"fmt"
	"net/http"
	"time"
)

const (
	GarantexApiUrl = "https://garantex.org/api/v2/depth?market=usdtrub"
)

type GarantexAPIResponse struct {
	Asks      []PriceInstance `json:"asks"`
	Bids      []PriceInstance `json:"bids"`
	Timestamp int             `json:"timestamp"`
}

type PriceInstance struct {
	Price string `json:"price"`
}

type GarantexFetcher struct {
	client *http.Client
}

func NewGarantexFetcher() *GarantexFetcher {

	c := &http.Client{
		Timeout: 5 * time.Second,
	}

	return &GarantexFetcher{
		client: c,
	}
}

func (r GarantexFetcher) FetchRate(ctx context.Context) (*domain.Rate, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, GarantexApiUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var apiResponse GarantexAPIResponse

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	resp.Body.Close()

	if len(apiResponse.Asks) == 0 || len(apiResponse.Bids) == 0 {
		return nil, errors.New("not enough data in garantex API response")
	}

	ask := apiResponse.Asks[0].Price
	bid := apiResponse.Bids[0].Price

	return &domain.Rate{
		Ask:       ask,
		Bid:       bid,
		Timestamp: time.Unix(int64(apiResponse.Timestamp), 0),
	}, nil
}
