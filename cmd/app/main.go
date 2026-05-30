package main

import (
	"fmt"
	"net/http"

	"alvintanoto.id/go-template/internal/adapters/auth"
	"alvintanoto.id/go-template/internal/adapters/handlers"
	"alvintanoto.id/go-template/internal/adapters/postgres"
	"alvintanoto.id/go-template/internal/adapters/redis"
	"alvintanoto.id/go-template/internal/application"
	"alvintanoto.id/go-template/internal/config"
	"github.com/go-chi/chi/v5"
	"go.uber.org/dig"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func BuildContainer() *dig.Container {
	c := dig.New()

	// import config
	c.Provide(config.NewConfig)

	// initialize logger
	c.Provide(func(cfg *config.Config) (*zap.Logger, error) {
		if cfg.AppEnv == "production" {
			return zap.NewProduction()
		}
		return zap.NewDevelopment()
	})

	// initialize db
	c.Provide(postgres.NewGormDatabase)

	// initialize redis
	c.Provide(redis.NewRedisClient)

	// setup hasher
	c.Provide(auth.NewHasher)

	// setup services
	c.Provide(application.NewAuthService)

	// setup repository
	c.Provide(postgres.NewUserRepository)

	// setup apps layer

	//routing
	c.Provide(handlers.NewZapMiddleware)
	c.Provide(handlers.NewAuthHandler)
	c.Provide(handlers.NewRouter)

	return c
}

func main() {
	container := BuildContainer()

	err := container.Invoke(func(db *gorm.DB) {})

	err = container.Invoke(func(cfg *config.Config, router *chi.Mux, logger *zap.Logger) error {
		defer logger.Sync()

		addr := fmt.Sprintf(":%s", cfg.Port)
		server := &http.Server{
			Addr:    addr,
			Handler: router,
		}

		fmt.Printf("Server starting on %s in [%s] mode...\n", addr, cfg.AppEnv)
		return server.ListenAndServe()
	})

	if err != nil {
		panic(err)
	}
}
