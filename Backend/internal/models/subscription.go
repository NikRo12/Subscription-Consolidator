package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Subscription struct {
	ID    int
	Name  string
	URL   string
	Price int
}

func (s *Subscription) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.URL, validation.Required, is.URL),
	)
}
