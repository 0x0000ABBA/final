package service

import (
	"context"
	"errors"
	"final/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockRateSaver is a mock implementation of the RateSaver interface.
type MockRateSaver struct {
	mock.Mock
}

func (m *MockRateSaver) SaveRate(ctx context.Context, rate *domain.Rate) error {
	args := m.Called(ctx, rate)
	return args.Error(0)
}

// MockRateFetcher is a mock implementation of the RateFetcher interface.
type MockRateFetcher struct {
	mock.Mock
}

func (m *MockRateFetcher) FetchRate(ctx context.Context) (*domain.Rate, error) {
	args := m.Called(ctx)
	if rate, ok := args.Get(0).(*domain.Rate); ok {
		return rate, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestRateService_GetRate(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	sugaredLogger := logger.Sugar()

	mockSaver := new(MockRateSaver)
	mockFetcher := new(MockRateFetcher)

	service := NewRateService(mockSaver, mockFetcher, sugaredLogger)

	fixedTime := time.Unix(1700000000, 0)
	rate := &domain.Rate{
		Ask:       "100.5",
		Bid:       "99.5",
		Timestamp: fixedTime,
	}

	tests := []struct {
		name          string
		fetcherMock   func()
		saverMock     func()
		expectedRate  *domain.Rate
		expectedError error
	}{
		{
			name: "Successful fetch and save",
			fetcherMock: func() {
				mockFetcher.On("FetchRate", mock.Anything).Return(rate, nil)
			},
			saverMock: func() {
				mockSaver.On("SaveRate", mock.Anything, rate).Return(nil)
			},
			expectedRate:  rate,
			expectedError: nil,
		},
		{
			name: "Fetcher returns error",
			fetcherMock: func() {
				mockFetcher.On("FetchRate", mock.Anything).Return((*domain.Rate)(nil), errors.New("fetch error"))
			},
			saverMock:     func() {},
			expectedRate:  nil,
			expectedError: errors.New("failed to fetch rate: fetch error"),
		},
		{
			name: "Saver returns error",
			fetcherMock: func() {
				mockFetcher.On("FetchRate", mock.Anything).Return(rate, nil)
			},
			saverMock: func() {
				mockSaver.On("SaveRate", mock.Anything, rate).Return(errors.New("save error"))
			},
			expectedRate:  rate,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFetcher.ExpectedCalls = nil
			mockSaver.ExpectedCalls = nil

			tt.fetcherMock()
			tt.saverMock()

			got, err := service.GetRate(context.Background())

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedRate, got)

			mockFetcher.AssertExpectations(t)
			mockSaver.AssertExpectations(t)
		})
	}
}
