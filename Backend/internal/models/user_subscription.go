package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type UserSubscription struct {
	UserID          int
	SubID           int
	Price           decimal.Decimal
	Period          Period
	NextPaymentDate time.Time
	IsActive        bool
}
