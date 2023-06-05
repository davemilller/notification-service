package subscriptions

import (
	"context"

	"go.uber.org/zap"
)

type Subscriber struct {
	Ctx          context.Context
	Subscription chan interface{}
	ID           string
}

type SubscriptionManager struct {
	Subs   map[string]*Subscriber
	Add    chan *Subscriber
	Remove chan string
}

func NewSubscriptionManager() *SubscriptionManager {
	sm := &SubscriptionManager{
		Subs:   make(map[string]*Subscriber),
		Add:    make(chan *Subscriber, 1),
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

func (sm *SubscriptionManager) AddSubscriber(ctx context.Context, userID string) error {
	if _, ok := sm.Subs[userID]; ok {
		zap.S().Infof("subscriber already exists: %s", userID)
		return nil
	}

	ch := make(chan interface{}, 1)

	sm.Add <- &Subscriber{
		Ctx:          ctx,
		Subscription: ch,
		ID:           userID,
	}

	go func() {
		<-ctx.Done()

		sm.Remove <- userID
	}()

	return nil
}

func (sm *SubscriptionManager) Push(userID string, note interface{}) error {
	if _, ok := sm.Subs[userID]; !ok {
		zap.S().Infof("subscriber does not exist: %s", userID)
		return nil
	}

	sm.Subs[userID].Subscription <- note

	return nil
}
