package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

func TestGigaChatClient_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping real API integration test in short mode.")
	}

	if err := godotenv.Load("../../../.env"); err != nil {
		t.Logf("Warning: Could not load .env file: %v. (This is fine if running in CI/CD with injected secrets)", err)
	}

	certPath := "../../../certs/russian_trusted_root_ca_pem.crt"

	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		t.Skipf("Certificate not found at %s. Skipping test.", certPath)
	}

	client, err := NewGigaChatClient(certPath)
	if err != nil {
		t.Fatalf("Failed to create GigaChat client: %v", err)
	}

	ctx := context.Background()

	emailText := "Thank you for your purchase! Your Spotify Premium subscription of $9.99 was charged on 2023-11-05."

	prompt := fmt.Sprintf("Extract the subscription information from this text: '%s'. "+
		"Return ONLY a valid JSON object with the following string fields: 'serviceName', 'price', 'date'. "+
		"Do not add any markdown formatting, explanations, or extra text.", emailText)

	t.Log("Sending request to real GigaChat API... (this might take a few seconds)")
	responseString, err := client.SendPrompt(ctx, prompt)
	if err != nil {
		t.Fatalf("SendPrompt returned an error: %v", err)
	}

	t.Logf("Raw GigaChat Response:\n%s", responseString)

	cleanJSON := strings.TrimSpace(responseString)
	cleanJSON = strings.TrimPrefix(cleanJSON, "```json")
	cleanJSON = strings.TrimPrefix(cleanJSON, "```")
	cleanJSON = strings.TrimSuffix(cleanJSON, "```")
	cleanJSON = strings.TrimSpace(cleanJSON)

	var resultData map[string]interface{}
	if err := json.Unmarshal([]byte(cleanJSON), &resultData); err != nil {
		t.Fatalf("Failed to parse GigaChat response as JSON: %v\nCleaned Response was: %s", err, cleanJSON)
	}

	expectedFields := []string{"serviceName", "price", "date"}
	for _, field := range expectedFields {
		if _, exists := resultData[field]; !exists {
			t.Errorf("Expected field %q is missing from the JSON response", field)
		}
	}

	if serviceName, ok := resultData["serviceName"].(string); ok {
		if !strings.Contains(strings.ToLower(serviceName), "spotify") {
			t.Errorf("Expected serviceName to contain 'spotify', but got: %s", serviceName)
		}
	} else {
		t.Errorf("serviceName was not a string")
	}
}
