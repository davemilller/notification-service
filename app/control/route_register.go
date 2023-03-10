package control

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/graphql-go/handler"
	"go.uber.org/zap"
)

func Routes(mc *Controller) []func(r *mux.Router) {
	return []func(r *mux.Router){func(r *mux.Router) {
		RegisterRoutes(mc, r)
	}, func(r *mux.Router) {
		RegisterNoAuthRoutes(mc, r)
	}}
}

func RegisterRoutes(mc *Controller, r *mux.Router) {
	zap.S().Info("registering authorized routes")

	// graphql setup
	gqlSchema, err := GenSchema(mc.gc)
	if err != nil {
		zap.S().Fatal("gqlschema: ", err)
	}

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

	// rest endpoints go here
}

func RegisterNoAuthRoutes(mc *Controller, r *mux.Router) {
	zap.S().Info("registering unauthorized routes")

	// add healthcheck
	r.Handle("/healthcheck", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "health check ok\n")
	})).Methods(http.MethodGet)
}
