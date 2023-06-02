package control

import (
	"context"

	"github.com/davemilller/notification-service/domain"
	"github.com/davemilller/notification-service/framework/repo"
)

var _ NotificationService = &repo.NotificationRepo{}

type NotificationService interface {
	Push(ctx context.Context, note *domain.Notification) error
	Get(ctx context.Context, userID string) ([]domain.Notification, error)
	Ack(ctx context.Context, id int) error
}
