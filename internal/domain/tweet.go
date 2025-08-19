package domain

import "errors"

type Tweet struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Text      string `json:"text"`
	CreatedAt int64  `json:"created_at"`
}

func NewTweet(id, userID, text string, createdAt int64) (Tweet, error) {
	if userID == "" {
		return Tweet{}, errors.New("user_id required")
	}
	if l := len(text); l == 0 || l > 280 {
		return Tweet{}, errors.New("text length must be 1..280")
	}
	return Tweet{ID: id, UserID: userID, Text: text, CreatedAt: createdAt}, nil
}
