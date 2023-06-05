package control

import "github.com/graphql-go/graphql"

var NotificationGraph = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "notification",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"userID": &graphql.Field{
				Type: graphql.String,
			},
			"details": &graphql.Field{
				Type: graphql.String,
			},
			"timestamp": &graphql.Field{
				Type: Time,
			},
		},
	},
)
