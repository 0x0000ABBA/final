package service

import (
	"bytes"
	"context"
	"errors"
	"final/internal/domain"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"
)

type mockRoundTripper struct {
	mockResponse *http.Response
	mockError    error
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.mockResponse, m.mockError
}

func TestNewGarantexFetcher(t *testing.T) {
	tests := []struct {
		name string
		want *GarantexFetcher
	}{
		{
			name: "Initialize GarantexFetcher with default HTTP client",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewGarantexFetcher()
			if got.client == nil {
				t.Errorf("NewGarantexFetcher() client = nil, want non-nil *http.Client")
			}
			if got.client.Timeout != 5*time.Second {
				t.Errorf("NewGarantexFetcher() client timeout = %v, want %v", got.client.Timeout, 5*time.Second)
			}
		})
	}
}

func TestGarantexFetcher_FetchRate(t *testing.T) {
	fixedUnixTimestamp := int64(1_700_000_000)
	fixedTime := time.Unix(fixedUnixTimestamp, 0)

	type fields struct {
		client *http.Client
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mockResp  *http.Response
		mockErr   error
		want      *domain.Rate
		wantErr   bool
		errorText string
	}{
		{
			name: "Successful fetch and parse",
			fields: fields{
				client: &http.Client{
					Transport: &mockRoundTripper{
						mockResponse: &http.Response{
							StatusCode: http.StatusOK,
							Body: io.NopCloser(bytes.NewBufferString(`{
								"asks": [{"price": "100.5"}],
								"bids": [{"price": "99.5"}],
								"timestamp": 1700000000
							}`)),
							Header: make(http.Header),
						},
						mockError: nil,
					},
				},
			},
			args: args{
				ctx: context.Background(),
			},
			want: &domain.Rate{
				Ask:       "100.5",
				Bid:       "99.5",
				Timestamp: fixedTime,
			},
			wantErr: false,
		},
		{
			name: "HTTP request execution failure",
			fields: fields{
				client: &http.Client{
					Transport: &mockRoundTripper{
						mockResponse: nil,
						mockError:    errors.New("network error"),
					},
				},
			},
			args: args{
				ctx: context.Background(),
			},
			want:      nil,
			wantErr:   true,
			errorText: "failed to do request",
		},
		{
			name: "Invalid JSON response",
			fields: fields{
				client: &http.Client{
					Transport: &mockRoundTripper{
						mockResponse: &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewBufferString(`invalid json`)),
							Header:     make(http.Header),
						},
						mockError: nil,
					},
				},
			},
			args: args{
				ctx: context.Background(),
			},
			want:      nil,
			wantErr:   true,
			errorText: "failed to decode response",
		},
		{
			name: "Failed to create HTTP request",
			fields: fields{
				client: &http.Client{
					Transport: &mockRoundTripper{
						mockResponse: &http.Response{
							StatusCode: http.StatusOK,
							Body: io.NopCloser(bytes.NewBufferString(`{
								"asks": [{"price": "100.5"}],
								"bids": [{"price": "99.5"}],
								"timestamp": 1700000000
							}`)),
							Header: make(http.Header),
						},
						mockError: nil,
					},
				},
			},
			args: args{
				ctx: nil, // error due to invalid context
			},
			want:      nil,
			wantErr:   true,
			errorText: "failed to create request",
		},
		{
			name: "Empty asks in response",
			fields: fields{
				client: &http.Client{
					Transport: &mockRoundTripper{
						mockResponse: &http.Response{
							StatusCode: http.StatusOK,
							Body: io.NopCloser(bytes.NewBufferString(`{
								"asks": [],
								"bids": [{"price": "99.5"}],
								"timestamp": 1700000000
							}`)),
							Header: make(http.Header),
						},
						mockError: nil,
					},
				},
			},
			args: args{
				ctx: context.Background(),
			},
			want:      nil,
			wantErr:   true,
			errorText: "not enough data in garantex API response",
		},
		{
			name: "Empty bids in response",
			fields: fields{
				client: &http.Client{
					Transport: &mockRoundTripper{
						mockResponse: &http.Response{
							StatusCode: http.StatusOK,
							Body: io.NopCloser(bytes.NewBufferString(`{
								"asks": [{"price": "100.5"}],
								"bids": [],
								"timestamp": 1700000000
							}`)),
							Header: make(http.Header),
						},
						mockError: nil,
					},
				},
			},
			args: args{
				ctx: context.Background(),
			},
			want:      nil,
			wantErr:   true,
			errorText: "not enough data in garantex API response",
		},
		{
			name: "Non-200 HTTP response",
			fields: fields{
				client: &http.Client{
					Transport: &mockRoundTripper{
						mockResponse: &http.Response{
							StatusCode: http.StatusInternalServerError,
							Body:       io.NopCloser(bytes.NewBufferString(`Internal Server Error`)),
							Header:     make(http.Header),
						},
						mockError: nil,
					},
				},
			},
			args: args{
				ctx: context.Background(),
			},
			want:    &domain.Rate{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fetcher GarantexFetcher
			if tt.fields.client != nil {
				fetcher = GarantexFetcher{
					client: tt.fields.client,
				}
			} else {
				fetcher = GarantexFetcher{
					client: nil,
				}
			}

			got, err := fetcher.FetchRate(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchRate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				if !errors.Is(err, errors.New(tt.errorText)) && !contains(err.Error(), tt.errorText) {
					t.Errorf("FetchRate() error = %v, expected to contain %v", err, tt.errorText)
				}
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchRate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to check if a substring exists within a string.
func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
