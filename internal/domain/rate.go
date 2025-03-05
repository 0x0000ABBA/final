package domain

import (
	"time"
)

type Rate struct {
	Ask       string
	Bid       string
	Timestamp time.Time
}
