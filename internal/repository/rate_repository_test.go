package repository

import (
	"context"
	"final/internal/domain"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"reflect"
	"testing"
	"time"
)

func TestNewRateRepository(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	type args struct {
		db *sqlx.DB
	}
	tests := []struct {
		name string
		args args
		want *RateRepository
	}{
		{
			name: "Create RateRepository with valid db",
			args: args{db: sqlxDB},
			want: &RateRepository{db: sqlxDB},
		},
		{
			name: "Create RateRepository with nil db",
			args: args{db: nil},
			want: &RateRepository{db: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRateRepository(tt.args.db)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRateRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRateRepository_SaveRate(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx  context.Context
		rate *domain.Rate
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "Successful record insertion",
			fields: fields{
				db: sqlxDB,
			},
			args: args{
				ctx: context.Background(),
				rate: &domain.Rate{
					Ask:       "100.5",
					Bid:       "99.5",
					Timestamp: time.Now(),
				},
			},
			mock: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO "Rate" \("ask", "bid", "timestamp"\) VALUES \(\$1, \$2, \$3\)`
				mock.ExpectExec(query).
					WithArgs(
						"100.5",
						"99.5",
						sqlmock.AnyArg(),
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "Error during query execution",
			fields: fields{
				db: sqlxDB,
			},
			args: args{
				ctx: context.Background(),
				rate: &domain.Rate{
					Ask:       "100.5",
					Bid:       "99.5",
					Timestamp: time.Now(),
				},
			},
			mock: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO "Rate" \("ask", "bid", "timestamp"\) VALUES \(\$1, \$2, \$3\)`
				mock.ExpectExec(query).
					WithArgs(
						"100.5",
						"99.5",
						sqlmock.AnyArg(),
					).
					WillReturnError(fmt.Errorf("db error"))
			},
			wantErr: true,
		},
		{
			name: "No expected queries executed",
			fields: fields{
				db: sqlxDB,
			},
			args: args{
				ctx: context.Background(),
				rate: &domain.Rate{
					Ask:       "100.5",
					Bid:       "99.5",
					Timestamp: time.Now(),
				},
			},
			mock:    func(mock sqlmock.Sqlmock) {},
			wantErr: true,
		},
		{
			name: "Error preparing the query",
			fields: fields{
				db: sqlxDB,
			},
			args: args{
				ctx: context.Background(),
				rate: &domain.Rate{
					Ask:       "100.5",
					Bid:       "99.5",
					Timestamp: time.Now(),
				},
			},
			mock: func(mock sqlmock.Sqlmock) {
				query := `INSERT INTO "Rate" \("ask", "bid", "timestamp"\) VALUES \(\$1, \$2, \$3\)`
				mock.ExpectExec(query).
					WillReturnError(fmt.Errorf("prepare error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock(mock)
			}

			r := &RateRepository{
				db: tt.fields.db,
			}
			if err := r.SaveRate(tt.args.ctx, tt.args.rate); (err != nil) != tt.wantErr {
				t.Errorf("SaveRate() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unmet expectations: %s", err)
			}
		})
	}
}
