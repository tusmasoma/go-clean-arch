package main

import (
	"context"
	"fmt"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
	"go.uber.org/dig"

	"github.com/tusmasoma/go-clean-arch/config"
	handler "github.com/tusmasoma/go-clean-arch/interfaces/handler/http"
	middleware "github.com/tusmasoma/go-clean-arch/interfaces/middleware/http"
	"github.com/tusmasoma/go-clean-arch/repository/mysql"
	"github.com/tusmasoma/go-clean-arch/usecase"

	_ "github.com/go-sql-driver/mysql"
)

func HTTPBuildContainer(ctx context.Context) (*dig.Container, error) {
	container := dig.New()

	if err := container.Provide(func() context.Context {
		return ctx
	}); err != nil {
		log.Error("Failed to provide context")
		return nil, err
	}

	providers := []interface{}{
		config.NewServerConfig,
		config.NewDBConfig,
		mysql.NewMySQLDB,
		mysql.NewTransactionRepository,
		mysql.NewTaskRepository,
		usecase.NewTaskUseCase,
		handler.NewTaskHandler,
		func(
			serverConfig *config.ServerConfig,
			taskHandler handler.TaskHandler,
		) *chi.Mux {
			r := chi.NewRouter()
			r.Use(cors.Handler(cors.Options{
				AllowedOrigins:     []string{"https://*", "http://*"},
				AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Origin"},
				ExposedHeaders:     []string{"Link", "Authorization"},
				AllowCredentials:   true,
				MaxAge:             serverConfig.PreflightCacheDurationSec,
				OptionsPassthrough: true,
			}))
			r.Use(middleware.Logging)

			r.Route("/api", func(r chi.Router) {
				r.Route("/task", func(r chi.Router) {
					r.Get("/get", taskHandler.GetTask)
					r.Get("/list", taskHandler.ListTasks)
					r.Post("/create", taskHandler.CreateTask)
					r.Put("/update", taskHandler.UpdateTask)
					r.Delete("/delete", taskHandler.DeleteTask)
				})
			})

			return r
		},
	}

	for _, provider := range providers {
		if err := container.Provide(provider); err != nil {
			log.Critical("Failed to provide dependency", log.Fstring("provider", fmt.Sprintf("%T", provider)))
			return nil, err
		}
	}

	log.Info("Container built successfully")
	return container, nil
}
