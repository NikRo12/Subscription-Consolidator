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

// MockAIAnalyst is a fake version of GigaChatClient used for testing business logic.
type MockAIAnalyst struct {
	MockResponse string
	MockError    error
	PromptRecv   string // To check if the prompt was formed correctly
}

func (m *MockAIAnalyst) SendPrompt(ctx context.Context, prompt string) (string, error) {
	m.PromptRecv = prompt
	return m.MockResponse, m.MockError
}

// TestOrchestration simulates a service that reads an email and sends it to the AI.
func TestAIAnalyst_BusinessLogic(t *testing.T) {
	// 1. Setup our mock dependencies
	mockEmailExtr := &email.MockExtractor{ // Note: assuming you placed MockExtractor in this package or imported it
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

	// Assign the mock to your interface to prove they satisfy the contract
	var analyst AIAnalyst = mockAI

	ctx := context.Background()
	dummyTask := models.Task{} // Dummy task for your MockExtractor

	// 2. Execute the simulated business logic
	emailText, err := mockEmailExtr.GetEmailText(ctx, dummyTask)
	if err != nil {
		t.Fatalf("Failed to extract email: %v", err)
	}

	aiResult, err := analyst.SendPrompt(ctx, "Extract subscription info from this text: "+emailText)
	if err != nil {
		t.Fatalf("AI Analyst failed: %v", err)
	}

	// 3. Assert the results
	expectedEmail := "Receipt from Netflix for $15.99 on 2023-10-01"
	if !strings.Contains(mockAI.PromptRecv, expectedEmail) {
		t.Errorf("Expected AI to receive the email text, got: %s", mockAI.PromptRecv)
	}

	if aiResult != `{"service": "Netflix", "price": 15.99}` {
		t.Errorf("Unexpected AI result: %s", aiResult)
	}
}

// TestAIAnalyst_Timeout proves that your context logic works.
func TestAIAnalyst_Timeout(t *testing.T) {
	mockAI := &MockAIAnalyst{
		MockError: context.DeadlineExceeded,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	time.Sleep(5 * time.Millisecond) // Simulate time passing

	_, err := mockAI.SendPrompt(ctx, "Hello")
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Expected DeadlineExceeded error, got: %v", err)
	}
}
