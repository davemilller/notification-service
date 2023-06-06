package control

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/graphql-go/handler"
	"go.uber.org/zap"

	"github.com/eientei/wsgraphql/v1"
	"github.com/eientei/wsgraphql/v1/compat/gorillaws"
	"github.com/gorilla/websocket"
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

	// subscriptions / websocket handler
	subHandler, err := wsgraphql.NewServer(
		*gqlSchema,
		wsgraphql.WithKeepalive(time.Second*30),
		wsgraphql.WithConnectTimeout(time.Second*30),
		wsgraphql.WithUpgrader(gorillaws.Wrap(&websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			Subprotocols: []string{
				wsgraphql.WebsocketSubprotocolGraphqlWS.String(),
				wsgraphql.WebsocketSubprotocolGraphqlTransportWS.String(),
			},
		})),
	)
	if err != nil {
		panic(err)
	}

	// graphql handler
	h := handler.New(&handler.Config{
		Schema: gqlSchema,
		Pretty: true,
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
