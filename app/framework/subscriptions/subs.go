package subscriptions

import (
	"context"

	"github.com/davemilller/notification-service/domain"
	"go.uber.org/zap"
)

type SubscriptionManager struct {
	Subs   map[string]*domain.Subscriber
	Add    chan *domain.Subscriber
	Remove chan string
}

func NewSubscriptionManager() *SubscriptionManager {
	sm := &SubscriptionManager{
		Subs:   make(map[string]*domain.Subscriber),
		Add:    make(chan *domain.Subscriber, 1),
		Remove: make(chan string, 1),
	}

	go func() {
		for {
			select {
			case sub := <-sm.Add:
				sm.Subs[sub.ID] = sub
				zap.S().Infof("added subscriber: %s", sub.ID)

			case id := <-sm.Remove:
				close(sm.Subs[id].Subscription)
				delete(sm.Subs, id)

				zap.S().Infof("removed subscriber: %s", id)
			}
		}
	}()

	return sm
}

func (sm *SubscriptionManager) AddSubscriber(ctx context.Context, userID string) (*domain.Subscriber, error) {
	if s, ok := sm.Subs[userID]; ok {
		zap.S().Infof("subscriber already exists: %s", userID)
		return s, nil
	}

	ch := make(chan interface{}, 1)

	sub := &domain.Subscriber{
		Ctx:          ctx,
		Subscription: ch,
		ID:           userID,
	}

	sm.Add <- sub

	go func() {
		<-ctx.Done()

		sm.Remove <- userID
	}()

	return sub, nil
}

func (sm *SubscriptionManager) Push(userID string, note domain.Notification) error {
	zap.S().Infof("pushing notification %+v to subscriber: %s", note, userID)
	if _, ok := sm.Subs[userID]; !ok {
		zap.S().Infof("subscriber does not exist: %s", userID)
		return nil
	}

	select {
	case <-sm.Subs[userID].Ctx.Done():
	case sm.Subs[userID].Subscription <- note:
	}

	return nil
}
