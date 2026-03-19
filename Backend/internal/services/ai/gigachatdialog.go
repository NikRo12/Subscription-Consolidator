package ai

// describes a message that a client send to GigaChat AI Model - https://developers.sber.ru/docs/ru/gigachat/api/reference/rest/post-chat
type gchatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type gchatGenRequest struct {
	Model       string         `json:"model"`
	Messages    []gchatMessage `json:"messages"`
	Temperature float64        `json:"temperature"`
}

type gchatResponse struct {
	Choices []struct {
		Message gchatMessage `json:"message"`
	} `json:"choices"`
}
