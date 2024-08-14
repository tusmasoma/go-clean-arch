package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
	"go.uber.org/dig"

	"github.com/tusmasoma/go-clean-arch/config"
	handler "github.com/tusmasoma/go-clean-arch/interfaces/handler/gin"
	middleware "github.com/tusmasoma/go-clean-arch/interfaces/middleware/gin"
	"github.com/tusmasoma/go-clean-arch/repository/mysql"
	"github.com/tusmasoma/go-clean-arch/usecase"

	_ "github.com/go-sql-driver/mysql"
)

func GinBuildContainer(ctx context.Context) (*dig.Container, error) {
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
		// This is database-agnostic and can be swapped with another database like PostgreSQL
		mysql.NewMySQLDB,
		mysql.NewTransactionRepository,
		mysql.NewTaskRepository,
		usecase.NewTaskUseCase,
		handler.NewTaskHandler,
		func(
			serverConfig *config.ServerConfig,
			taskHandler handler.TaskHandler,
		) *gin.Engine {
			r := gin.Default()

			r.Use(cors.New(cors.Config{
				AllowOrigins:     []string{"https://*", "http://*"},
				AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Origin"},
				ExposeHeaders:    []string{"Link", "Authorization"},
				AllowCredentials: true,
				MaxAge:           time.Duration(serverConfig.PreflightCacheDurationSec) * time.Second,
			}))

			r.Use(middleware.Logging())

			api := r.Group("/api")
			{
				task := api.Group("/task")
				{
					task.GET("/get", taskHandler.GetTask)
					task.GET("/list", taskHandler.ListTasks)
					task.POST("/create", taskHandler.CreateTask)
					task.PUT("/update", taskHandler.UpdateTask)
					task.DELETE("/delete", taskHandler.DeleteTask)
				}
			}

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
