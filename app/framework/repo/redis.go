package repo

import (
	"context"
	"fmt"
	"notification-service/domain"

	"github.com/go-redis/redis"
)

type NotificationRepo struct {
	redis *redis.Client
}

func NewNotificationRepo(client *redis.Client) (*NotificationRepo, error) {
	if client == nil {
		return nil, fmt.Errorf("redis client is nil")
	}

	return &NotificationRepo{
		redis: client,
	}, nil
}

func (r *NotificationRepo) Push(ctx context.Context, note *domain.Notification) (*domain.Notification, error) {
	return nil, fmt.Errorf("implement me")
}

func (r *NotificationRepo) Get(ctx context.Context, userID string) ([]*domain.Notification, error) {
	return nil, fmt.Errorf("implement me")
}

func (r *NotificationRepo) Ack(ctx context.Context, id int) error {
	return fmt.Errorf("implement me")
}
