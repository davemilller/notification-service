package control

import (
	"net/http"

	"github.com/functionalfoundry/graphqlws"
	"github.com/gorilla/mux"
	"github.com/graphql-go/handler"
	"go.uber.org/zap"
)

func Routes(mc *Controller) func(r *mux.Router) {
	return func(r *mux.Router) {
		RegisterRoutes(mc, r)
	}
}

func RegisterRoutes(mc *Controller, r *mux.Router) {
	zap.S().Info("registering routes")

	// graphql setup
	gqlSchema, err := GenSchema(mc.gc)
	if err != nil {
		zap.S().Fatal("gqlschema: ", err)
	}

	mc.gc.Subs = graphqlws.NewSubscriptionManager(gqlSchema)
	subHandler := graphqlws.NewHandler(graphqlws.HandlerConfig{
		SubscriptionManager: mc.gc.Subs,
		Authenticate: func(token string) (interface{}, error) {
			return token, nil
		},
	})

	h := handler.New(&handler.Config{
		Schema:     gqlSchema,
		Pretty:     true,
		Playground: true,
	})

	wrapHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	}

	r.HandleFunc("/graphql", wrapHandler).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	r.Handle("/subscriptions", subHandler)
}
