package email

import (
	"context"
	"fmt"
	"time"

	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/models"
)

type MockExtractor struct {
	MockReplySecTime time.Duration
	ServiceName      string
	ServicePrice     string
	Date             string
	Format           string
}

func (m *MockExtractor) GetEmailText(
	ctx context.Context,
	task models.Task,
) (string, error) {
	// имитация ответа от GMAIL API
	timer := time.NewTimer(m.MockReplySecTime * time.Second)
	defer timer.Stop()

	select {
	case <-timer.C:
		emailText := fmt.Sprintf(m.Format, m.ServiceName, m.ServicePrice, m.Date)
		return emailText, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
