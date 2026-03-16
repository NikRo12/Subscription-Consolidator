package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Entry struct {
	UserID          int             `json:"uuid"`
	Title           string          `json:"title"`
	Price           decimal.Decimal `json:"price"`
	Currency        string          `json:"currency"`
	Period          Period          `json:"period"`
	Category        Category        `json:"category"`
	NextPaymentDate time.Time       `json:"next_payment_date"`
	IconURL         string          `json:"icon_url"`
	BrandColor      string          `json:"brand_color"`
	IsActive        bool            `json:"is_active"`
	Description     string          `json:"description"`
}
