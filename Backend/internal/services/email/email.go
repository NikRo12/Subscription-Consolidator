package email

import (
	"context"

	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/models"
)

type EmailExtractor interface {
	// Extract the email-text by its-id
	GetEmailText(
		ctx context.Context,
		task models.Task,
	) (
		string,
		error,
	)
}
