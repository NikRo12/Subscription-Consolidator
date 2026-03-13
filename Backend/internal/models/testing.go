package models

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestUser(t *testing.T) *User {
	return &User{
		Email:        "user@example.org",
		RefreshToken: "123456789",
	}
}

func TestSubscription(t *testing.T) *Subscription {
	return &Subscription{
		Title:       "test",
		Currency:    "RUB",
		Category:    "entertainment",
		IconURL:     "https://logo.clearbit.com/netflix.com",
		BrandColor:  "#E50914",
		Description: "Тариф 'Standard', списание через 3 дня",
	}
}

func TestUserSubscription(t *testing.T) *UserSubscription {
	return &UserSubscription{
		UserID:          1,
		SubID:           1,
		Price:           decimal.NewFromFloat(899.99),
		Period:          "monthly",
		NextPaymentDate: time.Now(),
		IsActive:        true,
	}
}
