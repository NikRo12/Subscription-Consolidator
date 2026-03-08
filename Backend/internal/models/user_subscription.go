package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type UserSubscription struct {
	UserID          int
	SubID           int
	Period          Period
	Price           decimal.Decimal
	NextPaymentDate time.Time
	IsActive        bool
}
