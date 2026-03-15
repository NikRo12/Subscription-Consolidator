package ai

import (
	"context"
)

type AIAnalyst interface {
	SendPrompt(
		context.Context,
		string,
	) (response string,
		err error)
}
