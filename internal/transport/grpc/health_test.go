package grpc

import (
	"context"
	"final/internal/transport/gen"
	"reflect"
	"testing"
)

func TestHealthServiceServer_HealthCheck(t *testing.T) {
	type fields struct {
		UnimplementedHealthServiceServer gen.UnimplementedHealthServiceServer
	}
	type args struct {
		ctx context.Context
		req *gen.HealthCheckRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *gen.HealthCheckResponse
		wantErr bool
	}{
		{
			name:   "Any request should return valid response",
			fields: fields{},
			args: args{
				ctx: nil,
				req: nil,
			},
			want:    &gen.HealthCheckResponse{OK: true},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &HealthServiceServer{
				UnimplementedHealthServiceServer: tt.fields.UnimplementedHealthServiceServer,
			}
			got, err := s.HealthCheck(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("HealthCheck() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HealthCheck() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewHealthServiceServer(t *testing.T) {
	tests := []struct {
		name string
		want *HealthServiceServer
	}{
		{
			name: "Any call should return valid server",
			want: &HealthServiceServer{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHealthServiceServer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHealthServiceServer() = %v, want %v", got, tt.want)
			}
		})
	}
}
