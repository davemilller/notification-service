package domain

import "time"

type Notification struct {
	ID        int       `json:"id"`
	UserID    string    `json:"userID"`
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"timestamp"`
}
