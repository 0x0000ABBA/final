package grpc

import (
	"context"
	"errors"
	"final/internal/domain"
	"final/internal/transport/gen"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRateService struct {
	mock.Mock
}

func (m *MockRateService) GetRate(ctx context.Context) (*domain.Rate, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Rate), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestRateServiceServer_GetRate(t *testing.T) {
	mockService := new(MockRateService)
	server := NewRateServiceServer(mockService)

	fixedTime := time.Unix(1700000000, 0)

	rate := &domain.Rate{
		Ask:       "100.5",
		Bid:       "99.5",
		Timestamp: fixedTime,
	}

	tests := []struct {
		name          string
		setup         func()
		expectedResp  *gen.GetRateResponse
		expectedError string
	}{
		{
			name: "Successful GetRate",
			setup: func() {
				mockService.On("GetRate", mock.Anything).Return(rate, nil)
			},
			expectedResp: &gen.GetRateResponse{
				Ask:       rate.Ask,
				Bid:       rate.Bid,
				Timestamp: rate.Timestamp.String(),
			},
			expectedError: "",
		},
		{
			name: "RateService returns error",
			setup: func() {
				mockService.On("GetRate", mock.Anything).Return((*domain.Rate)(nil), errors.New("service error"))
			},
			expectedResp:  nil,
			expectedError: "error while using rate service: service error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.ExpectedCalls = nil
			tt.setup()

			req := &gen.GetRateRequest{}
			resp, err := server.GetRate(context.Background(), req)

			if tt.expectedError != "" {
				assert.Nil(t, resp)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResp, resp)
			}

			mockService.AssertExpectations(t)
		})
	}
}
func TestNewRateServiceServer(t *testing.T) {
	type args struct {
		service RateService
	}
	tests := []struct {
		name string
		args args
		want *RateServiceServer
	}{
		{
			name: "Nil service",
			args: args{
				service: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRateServiceServer(tt.args.service); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRateServiceServer() = %v, want %v", got, tt.want)
			}
		})
	}
}
