package parsercontroller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/NikRo12/Subscription-Consolidator/Backend/configs"
	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/models"
	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/services/ai"
	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/services/email"
	"github.com/NikRo12/Subscription-Consolidator/Backend/internal/services/filter"
	redisServ "github.com/NikRo12/Subscription-Consolidator/Backend/internal/services/redis"
	redisLib "github.com/redis/go-redis/v9"
)

func Run() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	redisConn := redisServ.NewRedisConn(redisLib.NewClient(
		&redisLib.Options{Addr: configs.GetRedisAddr()}),
	)
	defer redisConn.Terminate()

	certPath := configs.GetCertPath()

	analyzer, err := ai.NewGigaChatClient(certPath)
	if err != nil {
		log.Fatalf("cannot create a gigachat client: %v\n", err)
	}

	clientID, err := configs.GetGoogleClientID()
	if err != nil {
		log.Fatalf("cannot get clientID: %v\n", err)
	}

	clientSecret, err := configs.GetGoogleClientSecret()
	if err != nil {
		log.Fatalf("cannot get clientSecret: %v\n", err)
	}

	aiSemaphore := make(chan struct{}, 1)
	resChan := make(chan *models.ParseResult)

	var wgWorkers sync.WaitGroup

	go func() {
		for {
			task, err := redisConn.GetTask(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					log.Println("Stopping task fetcher...")
					return
				}
				log.Printf("cannot get task from redis %v\n", err)
				continue
			}

			wgWorkers.Add(1)
			go func(t *models.Task) {
				defer wgWorkers.Done()
				subroutine(
					t,
					analyzer,
					redisConn,
					resChan,
					clientID,
					clientSecret,
					aiSemaphore,
				)
			}(task)
		}
	}()

	go func() {
		<-ctx.Done()
		log.Println("Shutdown signal received. Waiting for active workers to complete...")

		wgWorkers.Wait()

		log.Println("All workers finished. Closing result channel...")
		close(resChan)
	}()

	log.Println("Parser Service is running. Press Ctrl+C to stop.")

	for res := range resChan {
		err := redisConn.PushParseResult(context.Background(), res)
		if err != nil {
			log.Printf("cannot push result: %v\n", err)
		} else {
			log.Printf("Successfully pushed result for UserID: %d\n", res.UserID)
		}
	}

	log.Println("Graceful shutdown complete. Bye!")
}

func subroutine(
	task *models.Task,
	analyzer *ai.GigaChatClient,
	redisConn *redisServ.RedisConn,
	resChan chan *models.ParseResult,
	clientID, clientSecret string,
	aiSemaphore chan struct{},
) {
	ctx := context.Background()

	gmailUser, err := email.ExtractGmailUser(ctx, task.RefreshToken, clientID, clientSecret)
	if err != nil {
		log.Printf("[UserID: %d] cannot extract gmail user: %v\n", task.UserID, err)
		return
	}

	rawEmails, err := gmailUser.GetEmailsText(20)
	if err != nil {
		log.Printf("[UserID: %d] cannot fetch emails: %v\n", task.UserID, err)
		return
	}

	filteredEmails := filter.FilterRelevantEmails(rawEmails)
	if len(filteredEmails) == 0 {
		log.Printf("[UserID: %d] no relevant emails found\n", task.UserID)
		return
	}

	var allEntries []models.Entry
	var mtx sync.Mutex
	var wg sync.WaitGroup

	for _, emailText := range filteredEmails {
		wg.Add(1)

		go func(text string) {
			defer wg.Done()

			aiSemaphore <- struct{}{}

			prompt := fmt.Sprintf(`Извлеки данные о подписке из этого чека.
Верни строго валидный JSON массив объектов. Каждый объект должен содержать поля:
- "title" (string): название сервиса/подписки
- "price" (number): сумма платежа
- "currency" (string): валюта (например "RUB", "USD")
- "period" (string): период оплаты (например "monthly", "yearly", "weekly")
- "category" (string): категория (например "streaming", "software", "music", "gaming")
- "next_payment_date" (string): дата следующего платежа в формате YYYY-MM-DD, если неизвестна — пустая строка ""
- "icon_url" (string): пустая строка ""
- "brand_color" (string): пустая строка ""
- "description" (string): краткое описание подписки
- "is_active" (bool): true
Если в тексте нет информации о подписке — верни []. Никаких пояснений, только JSON. Текст: %s`, text)

			ctxTmt, cancelTmt := context.WithTimeout(ctx, 15*time.Second)
			defer cancelTmt()

			aiResponse, err := analyzer.SendPrompt(ctxTmt, prompt)

			<-aiSemaphore

			if err != nil {
				log.Printf("[UserID: %d] AI analysis failed: %v\n", task.UserID, err)
				return
			}

			cleanJSON := extractJSONFromMarkdown(aiResponse)

			var entries []models.Entry
			if err := json.Unmarshal([]byte(cleanJSON), &entries); err != nil {
				log.Printf("[UserID: %d] Failed to unmarshal JSON from AI. Resp: %s\n", task.UserID, aiResponse)
				return
			}

			if len(entries) > 0 {
				mtx.Lock()
				allEntries = append(allEntries, entries...)
				mtx.Unlock()
			}
		}(emailText)
	}

	wg.Wait()

	if len(allEntries) > 0 {
		resChan <- &models.ParseResult{
			UserID:    task.UserID,
			EntryData: allEntries,
		}
	}
}

func extractJSONFromMarkdown(aiResponse string) string {
	resp := strings.TrimSpace(aiResponse)
	if strings.HasPrefix(resp, "```json") {
		resp = strings.TrimPrefix(resp, "```json")
		resp = strings.TrimSuffix(resp, "```")
	} else if strings.HasPrefix(resp, "```") {
		resp = strings.TrimPrefix(resp, "```")
		resp = strings.TrimSuffix(resp, "```")
	}
	return strings.TrimSpace(resp)
}
