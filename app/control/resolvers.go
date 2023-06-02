package control

import (
	"fmt"

	"github.com/davemilller/notification-service/domain"

	"github.com/graphql-go/graphql"
	"go.uber.org/zap"
)

func (gc *GQLController) HandleNotifications(p graphql.ResolveParams) (interface{}, error) {
	userID, ok := p.Args["userID"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid userID arg")
	}

	zap.S().Infof("hit subscription for userID: %s", userID)

	notes, err := gc.db.Get(p.Context, userID)
	if err != nil {
		return nil, err
	}

	subs := gc.Subs.Subscriptions()
	zap.S().Infof("subs: %+v", subs)
	for conn, sub := range subs {
		zap.S().Infof("conn: %+v", conn)
		zap.S().Infof("sub: %+v", sub)
	}

	return notes, nil
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
		UserID:  userID,
		Details: details,
	}

	err := gc.db.Push(p.Context, note)
	if err != nil {
		return nil, err
	}

	return note, nil
}
