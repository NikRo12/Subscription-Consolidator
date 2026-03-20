package email

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"sync"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type gmailUser struct {
	gmailService *gmail.Service
}

func ExtractGmailUser(
	ctx context.Context,
	refreshToken,
	clientID,
	clientSecret string,
) (*gmailUser, error) {

	if refreshToken == "" {
		return nil, fmt.Errorf("refresh token is empty (user needs to re-authenticate with consent)")
	}

	cfg := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
	}

	token := &oauth2.Token{RefreshToken: refreshToken}

	client := cfg.Client(ctx, token)

	gService, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("cannot create a gmail service: %w", err)
	}

	return &gmailUser{gmailService: gService}, nil
}

func (gu *gmailUser) getMessages(ctx context.Context) ([]*gmail.Message, error) {
	msgList, err := gu.gmailService.Users.Messages.List("me").Context(ctx).Do()

	if err != nil {
		return nil, fmt.Errorf("get user's messages list: %w", err)
	}

	return msgList.Messages, nil
}

/*
extract the text of last reqAmount messages
*/
func (gu *gmailUser) GetEmailsText(ctx context.Context, reqAmount int64) ([]string, error) {
	messages, err := gu.getMessages(ctx)

	if err != nil {
		return nil, fmt.Errorf("get messages: %w", err)
	}

	msgAmount := int64(len(messages))

	if msgAmount < reqAmount {
		reqAmount = msgAmount
	}

	extractedEmails := make([]string, 0, msgAmount)

	mtx := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	googleSemaphore := make(chan struct{}, 5)

	for _, msg := range messages[:reqAmount] {
		wg.Add(1)
		go func(msgID string) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				log.Printf("Context canceled before fetching msg %s: %v\n", msgID, ctx.Err())
				return
			case googleSemaphore <- struct{}{}:
			}

			defer func() { <-googleSemaphore }()

			fullMsg, err := gu.gmailService.Users.Messages.Get("me", msgID).Format("full").Context(ctx).Do()
			if err != nil {
				log.Printf("cannot fetch message %s: %v\n", msgID, err)
				return
			}

			bodyText := extractEmailBody(fullMsg.Payload)

			if bodyText == "" {
				bodyText = fullMsg.Snippet
			}

			mtx.Lock()
			extractedEmails = append(extractedEmails, bodyText)
			mtx.Unlock()

		}(msg.Id)
	}

	wg.Wait()

	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("fetching emails interrupted: %w", err)
	}

	return extractedEmails, nil
}

func extractEmailBody(part *gmail.MessagePart) string {
	if part == nil {
		return ""
	}

	if part.MimeType == "text/plain" && part.Body != nil && part.Body.Data != "" {
		decoded, err := base64.URLEncoding.DecodeString(part.Body.Data)
		if err != nil {
			return ""
		}
		return string(decoded)
	}

	for _, p := range part.Parts {
		text := extractEmailBody(p)
		if text != "" {
			return text
		}
	}

	return ""
}
