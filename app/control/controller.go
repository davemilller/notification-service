package control

type Controller struct {
	gc *GQLController
}

func NewController(db NotificationService) (*Controller, error) {
	graphqlController, err := NewGQLController(db)
	if err != nil {
		return nil, err
	}

	return &Controller{
		gc: graphqlController,
	}, nil
}
