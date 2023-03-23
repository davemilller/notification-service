package domain

import "time"

type Notification struct {
	ID        int       `json:"id"`
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"timestamp"`
}
