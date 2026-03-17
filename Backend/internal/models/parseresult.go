package models

type ParseResult struct {
	UserID    int     `json:"user-id"`
	EntryData []Entry `json:"entryData"`
}
