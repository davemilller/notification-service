package domain

import "context"

type Subscriber struct {
	Ctx          context.Context
	Subscription chan interface{}
	ID           string
}
