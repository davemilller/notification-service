package control

import (
	"fmt"

	"github.com/graphql-go/graphql"
)

func (gc *GQLController) HandleNotifications(p graphql.ResolveParams) (interface{}, error) {
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
