package models

type ParseResult struct {
	UserID    int     `json:"user_id"`
	EntryData []Entry `json:"entry_data"`
}
