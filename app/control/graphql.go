package control

import (
	"github.com/graphql-go/graphql"
	"go.uber.org/zap"
)

type GQLController struct {
	db   NotificationService
	subs SubscriberService
}

func NewGQLController(db NotificationService, subs SubscriberService) (*GQLController, error) {
	return &GQLController{
		db:   db,
		subs: subs,
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
			"getNotes": &graphql.Field{
				Type: graphql.NewList(NotificationGraph),
				Args: graphql.FieldConfigArgument{
					"userID": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: gc.GetNotes,
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
			"addNote": &graphql.Field{
				Type: NotificationGraph,
				Args: graphql.FieldConfigArgument{
					"userID": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"details": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: gc.AddNote,
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
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					zap.S().Infof("sub resolver")
					zap.S().Infof("p.Source: %+v", p.Source)
					return p.Source, nil
				},
				Subscribe: gc.AddSubscriber,
			},
		},
	}

	return graphql.NewObject(cfg)
}
