package control

import (
	"context"
	"notification-service/domain"
)

type NotificationService interface {
	Push(ctx context.Context, note *domain.Notification) (*domain.Notification, error)
	Get(ctx context.Context, userID string) ([]*domain.Notification, error)
	Ack(ctx context.Context, id int) error
}
