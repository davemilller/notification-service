package repo

import (
	"context"
	"fmt"

	"github.com/davemilller/notification-service/domain"

	"github.com/go-redis/redis"
)

type NotificationRepo struct {
	db *redis.Client
}

func NewNotificationRepo(client *redis.Client) (*NotificationRepo, error) {
	if client == nil {
		return nil, fmt.Errorf("redis client is nil")
	}

	return &NotificationRepo{
		db: client,
	}, nil
}

func (r *NotificationRepo) Push(ctx context.Context, note *domain.Notification) error {
	if note.UserID == "" {
		return fmt.Errorf("note is missing userID")
	}

	err := r.db.LPush(note.UserID, note)
	if err.Err() != nil {
		return err.Err()
	}

	return nil
}

func (r *NotificationRepo) Get(ctx context.Context, userID string) ([]domain.Notification, error) {
	if userID == "" {
		return nil, fmt.Errorf("missing userID")
	}

	var notes []domain.Notification
	err := r.db.Get(userID).Scan(&notes)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (r *NotificationRepo) Ack(ctx context.Context, id int) error {
	return fmt.Errorf("implement me")
}
