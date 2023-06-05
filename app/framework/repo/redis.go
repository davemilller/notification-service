package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/davemilller/notification-service/domain"
	"go.uber.org/zap"

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

	jsonNote, err := json.Marshal(note)
	if err != nil {
		return fmt.Errorf("marshalling note: %w", err)
	}

	redisErr := r.db.LPush(note.UserID, string(jsonNote))
	if redisErr.Err() != nil {
		return redisErr.Err()
	}

	return nil
}

func (r *NotificationRepo) Get(ctx context.Context, userID string) ([]domain.Notification, error) {
	if userID == "" {
		return nil, fmt.Errorf("missing userID")
	}

	notes, err := r.db.LRange(userID, 0, -1).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var notifications []domain.Notification
	for _, note := range notes {
		var n domain.Notification
		err := json.Unmarshal([]byte(note), &n)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	zap.S().Infof("got notes: %+v", notifications)

	return notifications, nil
}

func (r *NotificationRepo) Ack(ctx context.Context, id int) error {
	return fmt.Errorf("implement me")
}
