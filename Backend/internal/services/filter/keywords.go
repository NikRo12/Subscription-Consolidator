package filter

import (
	"strings"
)

var subscriptionKeywords = []string{
	"subscription", "receipt", "invoice", "payment", "renewal",
	"billing", "charged", "recurring", "order total", "transaction", "plan", "trial",
	"подписка", "чек", "оплата", "платеж", "списание",
	"тариф", "счет", "транзакция", "продление", "авторизация",
}

func BuildSearchQuery() string {
	keywordsJoined := strings.Join(subscriptionKeywords, " OR ")
	return "newer_than:1y AND (" + keywordsJoined + ")"
}
