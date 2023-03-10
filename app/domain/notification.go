package domain

import "time"

type Notification struct {
	ID        int
	Details   string
	CreatedAt time.Time
}
