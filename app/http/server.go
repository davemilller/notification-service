package http

import (
	"context"
	"errors"
	"net"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ListenAndServeDebug runs an HTTP debug pprof server.
func ListenAndServeDebug() error {
	return http.ListenAndServe(":6060", nil)
}

type Server struct {
	server          *http.Server
	router          *mux.Router
	port            string
	shutdownTimeout time.Duration
}

type ServerConfig struct {
	Port string        `env:"SERVER_PORT,default=:8080"`
	Wait time.Duration `env:"TYPE_DURATION,default=10s"`
}

func NewServerConfig() ServerConfig {
	return ServerConfig{}

}

func NewServer(ctx context.Context, cfg ServerConfig, addRoutes func(r *mux.Router)) (*Server, error) {
	s := &Server{
		server:          &http.Server{},
		router:          mux.NewRouter(),
		port:            cfg.Port,
		shutdownTimeout: cfg.Wait,
	}

	s.router.Handle("/favicon.ico", s.router.NotFoundHandler)

	// Our router is wrapped by another function handler to perform some
	// middleware-like tasks that cannot be performed by actual middleware.
	// This includes session handling
	s.server.Handler = http.HandlerFunc(s.serveHTTP)

	// add routes
	r := s.router.PathPrefix("/api/v1").Subrouter()

	// add middlewares
	addCorsMiddleware(r)

	addRoutes(r)

	return s, nil

}

func Invoke(lc fx.Lifecycle, s *Server) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			zap.S().Infof("Starting server on port %s", s.port)
			ln, err := net.Listen("tcp", s.port)
			if err != nil {
				return err
			}

			go func() {
				if err := s.server.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
					zap.S().Error(err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			zap.S().Info("Shutting server down...")
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			if err := s.server.Shutdown(ctx); err != nil {
				zap.S().Errorf("error shutting down server: %v", err)
			}
			zap.S().Info("server shutdown gracefully")

			return nil
		},
	})
}

func (s *Server) Open() (err error) {
	ln, err := net.Listen("tcp", s.port)
	if err != nil {
		return err
	}

	go func() {
		if err := s.server.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.S().Error(err)
		}
	}()

	return nil
}

// Close shuts down the server gracefully.
func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		zap.S().Errorf("error shutting down server: %v", err)
	} else {
		zap.S().Info("server shutdown gracefully")
	}

	return nil
}

func addCorsMiddleware(r *mux.Router) {
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(corsMiddleware(r))
}

func corsMiddleware(r *mux.Router) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "86400")
			if req.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, req)
		})
	}
}

func (s *Server) serveHTTP(w http.ResponseWriter, r *http.Request) {
	// Override method for forms passing "_method" value.
	if r.Method == http.MethodPost {
		switch v := r.PostFormValue("_method"); v {
		case http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete:
			r.Method = v
		}
	}

	// Override content-type for certain extensions.
	// This allows us to easily cURL API endpoints with a ".json" or ".csv"
	// extension instead of having to explicitly set Content-type & Accept headers.
	// The extensions are removed so they don't appear in the routes.
	switch ext := path.Ext(r.URL.Path); ext {
	case ".json":
		r.Header.Set("Accept", "application/json")
		r.Header.Set("Content-type", "application/json")
		r.URL.Path = strings.TrimSuffix(r.URL.Path, ext)
	case ".csv":
		r.Header.Set("Accept", "text/csv")
		r.URL.Path = strings.TrimSuffix(r.URL.Path, ext)
	}

	// Delegate remaining HTTP handling to the gorilla router.
	s.router.ServeHTTP(w, r)
}
