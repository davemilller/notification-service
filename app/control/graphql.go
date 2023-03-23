package control

import (
	"github.com/graphql-go/graphql"
)

type GQLController struct {
	db NotificationService
}

func NewGQLController(db NotificationService) (*GQLController, error) {
	return &GQLController{
		db: db,
	}, nil
}

func GenSchema(gc *GQLController) (*graphql.Schema, error) {
	// build schema
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:        gc.Query(),
		Mutation:     gc.Mutation(),
		Subscription: gc.Subscription(),
	})
	if err != nil {
		return nil, err
	}

	return &schema, nil
}

// Queries --------------------

func (gc *GQLController) Query() *graphql.Object {
	cfg := graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"healthcheck": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "check!", nil
				},
			},
		},
	}

	return graphql.NewObject(cfg)
}

// Mutations ------------------

func (gc *GQLController) Mutation() *graphql.Object {
	cfg := graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"healthcheck": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "check!", nil
				},
			},
		},
	}

	return graphql.NewObject(cfg)
}

// Subscriptions ---------------

func (gc *GQLController) Subscription() *graphql.Object {
	cfg := graphql.ObjectConfig{
		Name: "Subscription",
		Fields: graphql.Fields{
			"notifications": &graphql.Field{
				Type: NotificationGraph,
				Args: graphql.FieldConfigArgument{
					"userID": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: gc.HandleNotifications,
			},
		},
	}

	return graphql.NewObject(cfg)
}
