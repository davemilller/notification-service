package control

type Controller struct {
	gc *GQLController
}

func NewController(db NotificationService, subs SubscriberService) (*Controller, error) {
	graphqlController, err := NewGQLController(db, subs)
	if err != nil {
		return nil, err
	}

	return &Controller{
		gc: graphqlController,
	}, nil
}
