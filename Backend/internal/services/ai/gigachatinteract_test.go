package ai

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

// mockTransport implements http.RoundTripper. It intercepts all HTTP requests
// made by our http.Client and lets us define custom responses.
type mockTransport struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

func TestGigaChatClient_SendPrompt_Success(t *testing.T) {
	// 1. Create a custom RoundTripper to simulate the GigaChat API
	transport := &mockTransport{
		roundTripFunc: func(req *http.Request) (*http.Response, error) {
			// A. Simulate the OAuth Token response
			if req.URL.Path == "/api/v2/oauth" {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"access_token": "fake-jwt-token", "expires_at": 9999999999999}`)),
					Header:     make(http.Header),
				}, nil
			}

			// B. Simulate the Chat Completion response
			if req.URL.Path == "/api/v1/chat/completions" {
				// Verify the client correctly attached the Bearer token!
				authHeader := req.Header.Get("Authorization")
				if authHeader != "Bearer fake-jwt-token" {
					return &http.Response{
						StatusCode: http.StatusUnauthorized,
						Body:       io.NopCloser(bytes.NewBufferString(`{"error": "Unauthorized"}`)),
					}, nil
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(bytes.NewBufferString(`{
						"choices": [{
							"message": {
								"role": "assistant",
								"content": "I am GigaChat Lite, how can I help?"
							}
						}]
					}`)),
					Header: make(http.Header),
				}, nil
			}

			// Catch-all for unknown routes
			return &http.Response{StatusCode: http.StatusNotFound}, nil
		},
	}

	// 2. Build the client manually for the test to inject our mock transport
	// (This avoids having to read the real .env and russian.crt files during testing)
	client := &GigaChatClient{
		authKey:    "test-auth-key",
		httpClient: &http.Client{Transport: transport},
	}

	// 3. Run the actual method
	ctx := context.Background()
	response, err := client.SendPrompt(ctx, "Who are you?")

	// 4. Assert the result
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := "I am GigaChat Lite, how can I help?"
	if response != expected {
		t.Errorf("Expected %q, got %q", expected, response)
	}
}

func TestGigaChatClient_TokenCaching(t *testing.T) {
	oauthCallCount := 0

	transport := &mockTransport{
		roundTripFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path == "/api/v2/oauth" {
				oauthCallCount++
				// Return a token that is valid far into the future
				futureTime := time.Now().UnixMilli() + 3600000 // +1 hour
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"access_token": "cached-token", "expires_at": ` + fmt.Sprintf("%d", futureTime) + `}`)),
				}, nil
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"choices": [{"message": {"content": "ok"}}]}`)),
			}, nil
		},
	}

	client := &GigaChatClient{
		authKey:    "test-key",
		httpClient: &http.Client{Transport: transport},
	}

	ctx := context.Background()

	// Call SendPrompt 3 times
	client.SendPrompt(ctx, "Hello 1")
	client.SendPrompt(ctx, "Hello 2")
	client.SendPrompt(ctx, "Hello 3")

	// The token should have been fetched only ONCE because of your caching logic!
	if oauthCallCount != 1 {
		t.Errorf("Expected OAuth endpoint to be called exactly 1 time, but was called %d times", oauthCallCount)
	}
}
