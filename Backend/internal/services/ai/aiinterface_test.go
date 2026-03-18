package ai

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/models"
	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/services/email"
)

type MockAIAnalyst struct {
	MockResponse string
	MockError    error
	PromptRecv   string
}

func (m *MockAIAnalyst) SendPrompt(ctx context.Context, prompt string) (string, error) {
	m.PromptRecv = prompt
	return m.MockResponse, m.MockError
}

func TestAIAnalyst_BusinessLogic(t *testing.T) {
	mockEmailExtr := &email.MockExtractor{
		MockReplySecTime: 0,
		ServiceName:      "Netflix",
		ServicePrice:     "15.99",
		Date:             "2023-10-01",
		Format:           "Receipt from %s for $%s on %s",
	}

	mockAI := &MockAIAnalyst{
		MockResponse: `{"service": "Netflix", "price": 15.99}`,
		MockError:    nil,
	}

	var analyst AIAnalyst = mockAI

	ctx := context.Background()
	dummyTask := models.Task{}

	emailText, err := mockEmailExtr.GetEmailText(ctx, dummyTask)
	if err != nil {
		t.Fatalf("Failed to extract email: %v", err)
	}

	aiResult, err := analyst.SendPrompt(ctx, "Extract subscription info from this text: "+emailText)
	if err != nil {
		t.Fatalf("AI Analyst failed: %v", err)
	}

	expectedEmail := "Receipt from Netflix for $15.99 on 2023-10-01"
	if !strings.Contains(mockAI.PromptRecv, expectedEmail) {
		t.Errorf("Expected AI to receive the email text, got: %s", mockAI.PromptRecv)
	}

	if aiResult != `{"service": "Netflix", "price": 15.99}` {
		t.Errorf("Unexpected AI result: %s", aiResult)
	}
}

func TestAIAnalyst_Timeout(t *testing.T) {
	mockAI := &MockAIAnalyst{
		MockError: context.DeadlineExceeded,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	time.Sleep(5 * time.Millisecond)

	_, err := mockAI.SendPrompt(ctx, "Hello")
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Expected DeadlineExceeded error, got: %v", err)
	}
}
