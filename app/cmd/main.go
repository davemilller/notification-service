package main

import (
	"context"

	"github.com/davemilller/notification-service/control"
	"github.com/davemilller/notification-service/framework/repo"
	"github.com/davemilller/notification-service/framework/subscriptions"
	"github.com/davemilller/notification-service/http"

	env "github.com/Netflix/go-env"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	logger := control.NewLogger()
	zap.S().Info("Hello, I am notification-service")

	fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(
			LoadConfig,
			context.Background,
			repo.NewRedis,

			fx.Annotate(
				repo.NewNotificationRepo,
				fx.As(new(control.NotificationService)),
			),
			fx.Annotate(
				subscriptions.NewSubscriptionManager,
				fx.As(new(control.SubscriberService)),
			),

			control.NewController,
			control.Routes,
			http.NewServer,
			logger.Desugar,
		),
		fx.Invoke(
			http.Invoke,
		),
	).Run()
}

type AppConfig struct {
	Server http.ServerConfig
	DB     repo.DBConfig
}

func LoadConfig() (http.ServerConfig, repo.DBConfig, error) {
	cfg := AppConfig{
		Server: http.NewServerConfig(),
		DB:     repo.NewDBConfig(),
	}

	_, err := env.UnmarshalFromEnviron(&cfg)

	return cfg.Server, cfg.DB, err
}
