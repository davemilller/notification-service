package control

import "github.com/go-redis/redis"

type Controller struct {
	gc *GQLController
}

func NewController(db *redis.Client) (*Controller, error) {
	graphqlController, err := NewGQLController(db)
	if err != nil {
		return nil, err
	}

	return &Controller{
		gc: graphqlController,
	}, nil
}
