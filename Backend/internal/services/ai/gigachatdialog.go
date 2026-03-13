package ai

// describes a message that a client send to GigaChat AI Model - https://developers.sber.ru/docs/ru/gigachat/api/reference/rest/post-chat
type gchatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// describes structure of the body of an http-request to GigaChat API
type gchatGenRequest struct {
	Model       string         `json:"model"`
	Messages    []gchatMessage `json:"messages"`
	Temperature float64        `json:"temperature"`
}

// describes model's responses
type gchatResponse struct {
	Choices []struct {
		Message gchatMessage `json:"message"`
	} `json:"choices"`
}
