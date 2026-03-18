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

type mockTransport struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

func TestGigaChatClient_SendPrompt_Success(t *testing.T) {
	transport := &mockTransport{
		roundTripFunc: func(req *http.Request) (*http.Response, error) {
			if req.URL.Path == "/api/v2/oauth" {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"access_token": "fake-jwt-token", "expires_at": 9999999999999}`)),
					Header:     make(http.Header),
				}, nil
			}

			if req.URL.Path == "/api/v1/chat/completions" {
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

			return &http.Response{StatusCode: http.StatusNotFound}, nil
		},
	}

	client := &GigaChatClient{
		authKey:    "test-auth-key",
		httpClient: &http.Client{Transport: transport},
	}

	ctx := context.Background()
	response, err := client.SendPrompt(ctx, "Who are you?")

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
				futureTime := time.Now().UnixMilli() + 3600000
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

	client.SendPrompt(ctx, "Hello 1")
	client.SendPrompt(ctx, "Hello 2")
	client.SendPrompt(ctx, "Hello 3")

	if oauthCallCount != 1 {
		t.Errorf("Expected OAuth endpoint to be called exactly 1 time, but was called %d times", oauthCallCount)
	}
}
