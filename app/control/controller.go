package control

type Controller struct {
	gc *GQLController
}

func NewController() (*Controller, error) {
	graphqlController, err := NewGQLController()
	if err != nil {
		return nil, err
	}

	return &Controller{
		gc: graphqlController,
	}, nil
}
