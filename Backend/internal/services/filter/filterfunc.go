package filter

import (
	"strings"
	"sync"
)

func FilterRelevantEmails(emails []string) []string {
	var filtered []string
	var wg sync.WaitGroup
	var mtx sync.Mutex

	if len(emails) == 0 {
		return filtered
	}

	for _, emailText := range emails {
		wg.Add(1)

		go func(text string) {
			defer wg.Done()

			lowerText := strings.ToLower(text)

			for _, keyword := range subscriptionKeywords {
				if strings.Contains(lowerText, keyword) {
					mtx.Lock()
					filtered = append(filtered, text)
					mtx.Unlock()

					break
				}
			}
		}(emailText)
	}

	wg.Wait()

	return filtered
}
