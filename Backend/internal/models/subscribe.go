package models

import (
	"math/big"
)

type Subscribe struct {
	Tag         string
	Description string
	Price       big.Int
}
