package ai

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

var ErrParseCertError = errors.New("Cannot parse the certificate")

// Custom client for
type GigaChatClient struct {
	authKey         string
	httpClient      *http.Client
	lastAccessToken *gigaChatAccessToken
}

func NewGigaChatClient(russianCertPath string) (*GigaChatClient, error) {
	certData, err := os.ReadFile(russianCertPath)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Read certificate file: %w", err)
	}

	certPull := x509.NewCertPool()

	if ok := certPull.AppendCertsFromPEM(certData); !ok {
		return nil, ErrParseCertError
	}

	tlsClientConfig := &tls.Config{
		RootCAs: certPull,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsClientConfig,
	}

	authKey, err := getGigaChatAuthKey()
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Get auth key: %w", err)
	}

	return &GigaChatClient{authKey: authKey,
		httpClient:      &http.Client{Transport: transport, Timeout: time.Second * 30},
		lastAccessToken: nil}, nil
}

func (gch *GigaChatClient) updateAccessToken(ctx context.Context) error {
	if gch.lastAccessToken != nil &&
		gch.lastAccessToken.ExpiresAt > time.Now().UnixMilli()+time.Minute.Milliseconds() {
		return nil
	}

	url := "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"
	method := "POST"
	payload := strings.NewReader("scope=GIGACHAT_API_PERS")

	req, err := http.NewRequestWithContext(ctx, method, url, payload)
	if err != nil {
		return fmt.Errorf("[ERROR] NewRequest: %w", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("RqUID", uuid.New().String())
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", gch.authKey))

	res, err := gch.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do HTTP-Request to GigaChat API: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("API error: status code [%d]", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("[ERROR] ReadAll for http-Response body: %w", err)
	}

	accessToken := &gigaChatAccessToken{}
	if err := json.Unmarshal(body, accessToken); err != nil {
		return err
	}

	gch.lastAccessToken = accessToken

	return nil
}

func (gch *GigaChatClient) SendPrompt(ctx context.Context, prompt string) (string, error) {
	if err := gch.updateAccessToken(ctx); err != nil {
		return "", fmt.Errorf("failed to update token: %w", err)
	}

	url := "https://gigachat.devices.sberbank.ru/api/v1/chat/completions"
	method := "POST"

	body := &gchatGenRequest{
		Model: "GigaChat",
		Messages: []gchatMessage{
			{Role: "system", Content: "Ты - полезный AI ассистент, отвечающий строго по инструкциям пользователя"},
			{Role: "user", Content: prompt},
		},
		Temperature: 0.001,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create http request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", gch.lastAccessToken.AccessToken))

	res, err := gch.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		errorBody, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("API returned non-200 status [%d]: %s", res.StatusCode, string(errorBody))
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	chatResponse := &gchatResponse{}
	if err := json.Unmarshal(resBody, chatResponse); err != nil {
		return "", fmt.Errorf("failed to unmarshal API response: %w", err)
	}

	if len(chatResponse.Choices) == 0 {
		return "", errors.New("API returned successfully, but choices array is empty")
	}

	return chatResponse.Choices[0].Message.Content, nil
}
