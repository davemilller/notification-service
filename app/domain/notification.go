package domain

import (
	"encoding/json"
	"time"
)

type Notifications map[int][]Notification

type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userID"`
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"timestamp"`
}

func (n Notification) MarshalBinary() ([]byte, error) {
	return json.Marshal(n)
}
