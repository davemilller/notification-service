package control

import (
	"context"

	"github.com/davemilller/notification-service/domain"
	"github.com/davemilller/notification-service/framework/repo"
	"github.com/davemilller/notification-service/framework/subscriptions"
)

var _ NotificationService = &repo.NotificationRepo{}

type NotificationService interface {
	Push(ctx context.Context, note *domain.Notification) error
	Get(ctx context.Context, userID string) ([]domain.Notification, error)
	Ack(ctx context.Context, id int) error
}

var _ SubscriberService = &subscriptions.SubscriptionManager{}

type SubscriberService interface {
	AddSubscriber(ctx context.Context, userID string) error
	Push(userID string, note interface{}) error
}
