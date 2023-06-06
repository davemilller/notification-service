package control

import (
	"fmt"
	"time"

	"github.com/davemilller/notification-service/domain"
	"github.com/google/uuid"

	"github.com/graphql-go/graphql"
	"go.uber.org/zap"
)

func (gc *GQLController) AddSubscriber(p graphql.ResolveParams) (interface{}, error) {
	zap.S().Infof("add sub resolver")
	userID, ok := p.Args["userID"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid userID arg")
	}

	zap.S().Infof("hit subscription for userID: %s", userID)

	// add subscriber
	sub, err := gc.subs.AddSubscriber(p.Context, userID)
	if err != nil {
		return nil, err
	}

	return sub.Subscription, nil
}

func (gc *GQLController) AddNote(p graphql.ResolveParams) (interface{}, error) {
	userID, ok := p.Args["userID"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid userID arg")
	}

	details, ok := p.Args["details"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid details arg")
	}

	note := &domain.Notification{
		ID:        uuid.New().String(),
		UserID:    userID,
		Details:   details,
		CreatedAt: time.Now().Round(time.Microsecond).UTC(),
	}

	err := gc.db.Push(p.Context, note)
	if err != nil {
		return nil, err
	}

	err = gc.subs.Push(userID, *note)
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (gc *GQLController) GetNotes(p graphql.ResolveParams) (interface{}, error) {
	userID, ok := p.Args["userID"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid userID arg")
	}

	notes, err := gc.db.Get(p.Context, userID)
	if err != nil {
		return nil, err
	}

	return notes, nil
}
